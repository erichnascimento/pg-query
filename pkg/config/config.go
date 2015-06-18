package config

import(
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Host struct {
	Name string
	User string
	Password string 
	Host string
	Port uint16
	Connection string
	Databases []string
}

func (h *Host) DatabaseExists(name string) bool {
	for _, dbname := range h.Databases {
    	if dbname == name {
    		return true
    	}
  	}

  	return false
}

func (h *Host) ApplyDatabaseFilter(databases []string) error {
	if len(databases) == 0 {
		return nil
	}

	filteredDbs := make([]string, len(databases))
	for i, dbname := range databases {
		if !h.DatabaseExists(dbname) {
			return fmt.Errorf(`database not found "%s"`, dbname)
		}
		filteredDbs[i] = dbname
	}

	h.Databases = filteredDbs

	return nil
}

type Config struct {
	Hosts []*Host
}

func CreateFromFile(filePath string) (error, *Config) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %s", err), nil
	}

	// Unmarshal config
	config := new(Config)
	
	if err = yaml.Unmarshal(b, config); err != nil {
		return err, nil
	}

	return nil, config
}

func (c *Config) GetHostByName(name string) *Host {
	for _, h := range c.Hosts {
    	if h.Name == name {
    		return h
    	}
  	}

  	return nil
}

func (c *Config) ApplyHostFilter(hostNames []string) error {
	if len(hostNames) == 0 {
		return nil
	}

	filteredHosts := make([]*Host, len(hostNames))
	for i, name := range hostNames {
		h := c.GetHostByName(name)
		if h == nil {
			return fmt.Errorf(`host not found "%s"`, name)
		}
		filteredHosts[i] = h
	}

	c.Hosts = filteredHosts

	return nil
}