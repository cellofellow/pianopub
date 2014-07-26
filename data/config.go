package data

import (
	"encoding/json"

	"github.com/coopernurse/gorp"
)

type Config struct {
	Key       string                 `db:"key"`
	JSONField string                 `db:"json"`
	JSON      map[string]string      `db:"-"`
	dbmap     *gorp.DbMap            `db:"-"`
}


func (c Config) ddl(dbmap *gorp.DbMap) error {
	table := dbmap.AddTableWithName(Config{}, "config")
	table.SetKeys(false, "Key")
	table.ColMap("Key").SetNotNull(true)
	table.ColMap("JSONField").SetNotNull(true)

	_, err := dbmap.Exec(`
		CREATE TABLE IF NOT EXISTS config (
			key  TEXT NOT NULL PRIMARY KEY,
			json TEXT NOT NULL
		)
	`)

	return err
}

func (d *Database) AddConfig(key string, config map[string]string) (*Config, error) {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	c := &Config{Key: key, JSONField: string(jsonBytes)}
	err = d.dbmap.Insert(c)
	if err != nil {
		return nil, err
	}
	c.JSON = config
	c.dbmap = d.dbmap
	return c, nil
}

func (d *Database) GetConfig(key string) (*Config, error) {
	var c Config
	err := d.dbmap.SelectOne(&c, `SELECT * FROM config WHERE key = ?`, key)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(c.JSONField), c.JSON)
	if err != nil {
		return nil, err
	}
	c.dbmap = d.dbmap
	return &c, nil
}

func (c *Config) UpdateJSON(config map[string]string) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	jsonString := string(jsonBytes)
	_, err = c.dbmap.Exec(`UPDATE config SET json = ? WHERE key = ?`, jsonString, c.Key)
	if err != nil {
		return err
	}
	c.JSONField = jsonString
	c.JSON = config
	return nil
}
