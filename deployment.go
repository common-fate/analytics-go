package analytics

import (
	"context"
	"sync"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
)

var (
	depMutex     = &sync.RWMutex{}
	globalDeploy *Deployment
)

type globalDep struct {
	mu     *sync.Mutex
	dep    *Deployment
	loader DeploymentLoader
}

type DeploymentLoader interface {
	LoadDeployment(ctx context.Context) (*Deployment, error)
}

func (gd *globalDep) SetLoader(loader DeploymentLoader) {
	gd.loader = loader
}

func (gd *globalDep) Get(ctx context.Context) *Deployment {
	// if we don't have a deployment and there is no loader defined,
	// we can't load one. So just return with nil
	if gd.dep == nil && gd.loader == nil {
		return nil
	}

	// if we don't have a deployment, but we have a loader defined to fetch one,
	// try and fetch it.
	if gd.dep == nil {
		gd.loadDeployment(ctx)
	}

	// if the deployment ID is empty also return nil.
	if gd.dep.ID == "" {
		return nil
	}

	// if we get here, we've got a good deployment.
	return gd.dep
}

func (gd *globalDep) loadDeployment(ctx context.Context) {
	if gd.loader == nil {
		return
	}

	log := zap.L().Named("cf-analytics")
	d, err := gd.loader.LoadDeployment(ctx)
	if err != nil {
		log.Error("error loading deployment", zap.Error(err))
		// just set an empty deployment to prevent trying to load it again the next time.
		d = &Deployment{}
	}

	gd.mu.Lock()
	defer gd.mu.Unlock()

	gd.dep = d
}

var (
	globalDeployment = &globalDep{
		mu: &sync.Mutex{},
	}
)

// Deployment is a Common Fate deployment identifier.
// If you're editing this make sure you edit the Traits()
// method to ensure the properties propagate.
type Deployment struct {
	ID         string
	Version    string
	UserCount  int
	GroupCount int
}

// Traits returns the traits to use for the group identifier
func (d Deployment) Traits() acore.Traits {
	return acore.NewTraits().
		Set("version", d.Version).
		Set("user_count", d.UserCount).
		Set("group_count", d.GroupCount).
		Set("groupType", "deployment").
		Set("id", d.ID)
}

func getDeployment() *Deployment {
	depMutex.RLock()
	defer depMutex.RUnlock()
	d := globalDeploy
	return d
}

// SetDeployment sets deployment information.
func SetDeployment(dep *Deployment) {
	depMutex.Lock()
	defer depMutex.Unlock()
	globalDeploy = dep
}
