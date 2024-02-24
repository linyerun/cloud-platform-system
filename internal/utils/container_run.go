package utils

import "fmt"

type ContainerRunCommandOption func(options *containerRunCommandOptions)

type containerRunCommandOptions struct {
	commands         []string
	image            string
	containerCommand string
}

func CreateContainerRunCommand(options ...ContainerRunCommandOption) []string {
	obj := &containerRunCommandOptions{commands: []string{"run", "--privileged=true", "-itd"}}
	for _, option := range options {
		option(obj)
	}
	if len(obj.containerCommand) == 0 {
		return append(obj.commands, obj.image)
	}
	return append(obj.commands, obj.image, obj.containerCommand)
}

func WithPortMappingOption(from, to int64) ContainerRunCommandOption {
	return func(options *containerRunCommandOptions) {
		options.commands = append(options.commands, "-p", fmt.Sprintf("%d:%d", from, to))
	}
}

func WithNameOption(name string) ContainerRunCommandOption {
	return func(options *containerRunCommandOptions) {
		options.commands = append(options.commands, "--name", name)
	}
}

func WithCpuCoreCountOption(coreCnt uint) ContainerRunCommandOption {
	return func(options *containerRunCommandOptions) {
		options.commands = append(options.commands, "--cpus", fmt.Sprintf("%d", coreCnt))
	}
}

func WithMemoryOption(memory int64) ContainerRunCommandOption {
	return func(options *containerRunCommandOptions) {
		memory /= 1024
		options.commands = append(options.commands, "-m", fmt.Sprintf("%dM", memory))
	}
}

func WithMemorySwapOption(memorySwap int64) ContainerRunCommandOption {
	return func(options *containerRunCommandOptions) {
		if memorySwap >= 1024 {
			memorySwap /= 1024
		} else if memorySwap != -1 {
			memorySwap = 100
		} else {
			options.commands = append(options.commands, fmt.Sprintf("--memory-swap=%d", memorySwap))
			return
		}

		options.commands = append(options.commands, fmt.Sprintf("--memory-swap=%dM", memorySwap))
	}
}

func WithDiskSizeOption(diskSize uint) ContainerRunCommandOption {
	return func(options *containerRunCommandOptions) {
		// TODO 调试不行，后面可以了再搞
		//diskSize /= 1024
		//options.commands = append(options.commands, "--storage-opt", fmt.Sprintf("size=%dM", diskSize))
	}
}

func WithImageAndContainerCommand(image string, commands []string) ContainerRunCommandOption {
	return func(options *containerRunCommandOptions) {
		options.image = image
		for i, command := range commands {
			if i == 0 {
				options.containerCommand = fmt.Sprintf("%s", command)
				continue
			}
			options.containerCommand += fmt.Sprintf("&& %s", command)
		}
	}
}
