package modules

type Interface struct {
	IfaceName string `mapstructure:"interface_name"`
}

type Temperature struct {
	CpuTempPath string `mapstructure:"cpu_temp_path"`
}

type Config struct {
	Interface   Interface   `mapstructure:"interface"`
	Temperature Temperature `mapstructure:"temperature"`
}

func (c Interface) Configure() {
	ifaceName = c.IfaceName
}

func (c Temperature) Configure() {
	sysDirName = c.CpuTempPath
}

func Configure(c Config) {
	c.Interface.Configure()
	c.Temperature.Configure()
}
