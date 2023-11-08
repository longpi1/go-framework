package plugin

import (
"io/fs"
"os"
"path/filepath"

"go-pulgin/internal/infra"
"go-pulgin/internal/infra/config"
)

var system *infra.System

func StartSystem(pipelineConfPath string) {
	system = infra.NewSystem(config.NewYamlFactory())
	_ = filepath.WalkDir(pipelineConfPath, func(path string, d fs.DirEntry, err error) error {

		// 如果不是 YAML 文件，则跳过
		extension := filepath.Ext(path)
		if extension != ".yaml" && extension != ".yml" {
			return nil
		}

		// 读取 YAML 文件并解析
		fileContents, err := os.ReadFile(path)
		if err != nil {
			logger.Errorf("failed to read file: %v, skip it!", path)
			return err
		}

		err = system.LoadConf(string(fileContents))
		if err != nil {
			logger.Errorf("failed to load conf: %v, skip it!", path)
			return err
		}

		logger.Infof("success to load pipeline conf: %v", path)
		return nil
	})

	// 注册开关
	infra.NewSwitch(system).StartSwitchMonitor()
}
