package dialects

import (
	"fmt"
	"os"

	"github.com/arsura/goose/database/dialect"
)

const (
	DefaultReplicatedMergeTreeEngine = "ReplicatedMergeTree('/clickhouse/tables/{shard}/{database}/{table}', '{replica}')"
	DefaultMergeTreeEngine           = "MergeTree()"
)

// NewClickhouse returns a new [dialect.Querier] for Clickhouse dialect.
func NewClickhouse() dialect.Querier {
	var (
		cluster      = os.Getenv("CLICKHOUSE_CLUSTER")
		zooPath      = os.Getenv("CLICKHOUSE_ZOOKEEPER_PATH")
		replica      = os.Getenv("CLICKHOUSE_REPLICA_NAME")
		insertQuorum = os.Getenv("CLICKHOUSE_INSERT_QUORUM")
	)
	if insertQuorum == "" {
		insertQuorum = "auto"
	}
	return &clickhouse{cluster: cluster, zooPath: zooPath, replica: replica, insertQuorum: insertQuorum}
}

type clickhouse struct {
	cluster      string
	zooPath      string
	replica      string
	insertQuorum string
}

var _ dialect.Querier = (*clickhouse)(nil)

func (c *clickhouse) CreateTable(tableName string) string {
	cmd := c.getClusterCommand(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", tableName))
	engine := c.getTableEngine()
	q := `%s 
	(
		version_id Int64,
		is_applied UInt8,
		date Date default now(),
		tstamp DateTime default now()
	)
	ENGINE = %s
	ORDER BY (date)`
	return fmt.Sprintf(q, cmd, engine)
}

func (c *clickhouse) getClusterCommand(baseCommand string) string {
	if c.cluster != "" {
		return fmt.Sprintf("%s ON CLUSTER %s", baseCommand, c.cluster)
	}
	return baseCommand
}

func (c *clickhouse) getTableEngine() string {
	if c.cluster != "" {
		if c.zooPath != "" && c.replica != "" {
			return fmt.Sprintf("ReplicatedMergeTree('%s', '%s')", c.zooPath, c.replica)
		}
		return DefaultReplicatedMergeTreeEngine
	}
	return DefaultMergeTreeEngine
}

func (c *clickhouse) InsertVersion(tableName string) string {
	if c.cluster != "" {
		q := `INSERT INTO %s (version_id, is_applied) 
		SETTINGS insert_quorum=%s, insert_quorum_parallel=0, select_sequential_consistency=1
		VALUES ($1, $2)`
		return fmt.Sprintf(q, tableName, c.insertQuorum)
	}
	return fmt.Sprintf("INSERT INTO %s (version_id, is_applied) VALUES ($1, $2)", tableName)
}

func (c *clickhouse) DeleteVersion(tableName string) string {
	cmd := c.getClusterCommand(fmt.Sprintf("ALTER TABLE %s", tableName))
	q := `%s DELETE WHERE version_id = $1 SETTINGS mutations_sync = 2`
	return fmt.Sprintf(q, cmd)
}

func (c *clickhouse) GetMigrationByVersion(tableName string) string {
	q := `SELECT tstamp, is_applied FROM %s WHERE version_id = $1 ORDER BY tstamp DESC LIMIT 1`
	return fmt.Sprintf(q, tableName)
}

func (c *clickhouse) ListMigrations(tableName string) string {
	q := `SELECT version_id, is_applied FROM %s ORDER BY version_id DESC`
	return fmt.Sprintf(q, tableName)
}

func (c *clickhouse) GetLatestVersion(tableName string) string {
	q := `SELECT max(version_id) FROM %s`
	return fmt.Sprintf(q, tableName)
}
