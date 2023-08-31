package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	Project struct {
		Port           string `yaml:"port" env-default:"8001"`
		ProjectBaseDir string `yaml:"project_base_dir" env-default:"AvitoInternship"`
	} `yaml:"project"`

	Logger struct {
		LogsDir        string `yaml:"logs_dir" env-default:"logs/app/"`
		LogsFileName   string `yaml:"logs_file_name"`
		LogsUseStdOut  *bool  `yaml:"logs_use_std_out" env-default:"true"`
		LogsTimeFormat string `yaml:"logs_time_format" env-default:"2006-01-02_15:04:05_MST"`
	} `yaml:"logger"`

	DB struct {
		DBUser             string `env:"POSTGRES_USER"`
		DBPassword         string `env:"POSTGRES_PASSWORD"`
		DBHost             string `env:"POSTGRES_HOST"`
		DBPort             string `env:"POSTGRES_PORT"`
		DBSchemaName       string `env:"POSTGRES_SCHEMA"`
		DBUserTableName    string `yaml:"user_table_name" env-default:"users"`
		DBSegmentTableName string `yaml:"segment_table_name" env-default:"segments"`
		DBU2STableName     string `yaml:"u2s_table_name" env-default:"users2segments"`
		DBHistoryTableName string `yaml:"history_table_name" env-default:"history"`
		//DBTimeFormat       string `yaml:"time_format" env-default:"2006-01-02T15:04:05Z"`
	} `yaml:"db"`

	Routes struct {
		RoutePrefix string `yaml:"route_prefix"`

		// UserRoutes
		RouteUserCreate       string `yaml:"route_user_create" env-default:"/user/create"`
		RouteUser             string `yaml:"route_user" env-default:"/user/{id:[0-9]+}"`
		RouteUserSegments     string `yaml:"route_user_segments" env-default:"/user/{id:[0-9]+}/segments"`
		RouteUserEditSegments string `yaml:"route_user_edit_segments" env-default:"/user/{id:[0-9]+}/segments/edit"`

		// SegmentRoutes
		RouteSegmentCreate string `yaml:"route_segment_create" env-default:"/segment/create"`
		RouteSegment       string `yaml:"route_segment" env-default:"/segment/{slug}"`

		// History
		RouteGetHistory string `yaml:"route_history" env-default:"/history"`
	} `yaml:"routes"`

	//History struct {
	//	TimeFormat string `yaml:"time_format" env-default:"2006-01-02T15:04:05.999999Z"`
	//}
}

func Parse(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, errors.Wrap(err, "parse config")
	}

	return &cfg, nil
}
