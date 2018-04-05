package workflow

// NewDestroyWorkflow creates new instances of the 'destroy' workflow,
// responsible for running the actions required to remove resources
// of an existing cluster and clean up any remaining artefacts.
func NewDestroyWorkflow(clusterDir string) Workflow {
	return Workflow{
		metadata: metadata{clusterDir: clusterDir},
		steps: []Step{
			readClusterConfigStep,
			destroyJoinStep,
			destroyBootstrapStep,
			destroyAssetsStep,
		},
	}
}

func destroyAssetsStep(m *metadata) error {
	return runDestroyStep(m.clusterDir, assetsStep)
}

func destroyBootstrapStep(m *metadata) error {
	if err := runDestroyStep(m.clusterDir, etcdStep); err != nil {
		return err
	}

	if err := runDestroyStep(m.clusterDir, bootstrapStep); err != nil {
		return err
	}

	if err := destroyTNCDNS(m.clusterDir); err != nil {
		return err
	}

	return runDestroyStep(m.clusterDir, topologyStep)
}

func destroyJoinStep(m *metadata) error {
	if err := runDestroyStep(m.clusterDir, joinWorkersStep); err != nil {
		return err
	}
	return runDestroyStep(m.clusterDir, joinMastersStep)
}

func runDestroyStep(clusterDir, step string, extraArgs ...string) error {
	if !hasStateFile(clusterDir, step) {
		// there is no statefile, therefore nothing to destroy for this step
		return nil
	}
	templateDir, err := findTemplates(step)
	if err != nil {
		return err
	}

	return tfDestroy(clusterDir, step, templateDir, extraArgs...)
}
