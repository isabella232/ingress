package controller

import (
	"github.com/caddyserver/ingress/pkg/k8s"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

// ConfigMapAddedAction provides an implementation of the action interface.
type ConfigMapAddedAction struct {
	resource *v1.ConfigMap
}

// ConfigMapUpdatedAction provides an implementation of the action interface.
type ConfigMapUpdatedAction struct {
	resource    *v1.ConfigMap
	oldResource *v1.ConfigMap
}

// ConfigMapDeletedAction provides an implementation of the action interface.
type ConfigMapDeletedAction struct {
	resource *v1.ConfigMap
}

// onConfigMapAdded runs when a configmap is added to the namespace.
func (c *CaddyController) onConfigMapAdded(obj *v1.ConfigMap) {
	c.syncQueue.Add(ConfigMapAddedAction{
		resource: obj,
	})
}

// onConfigMapUpdated is run when a configmap is updated in the namespace.
func (c *CaddyController) onConfigMapUpdated(old *v1.ConfigMap, new *v1.ConfigMap) {
	c.syncQueue.Add(ConfigMapUpdatedAction{
		resource:    new,
		oldResource: old,
	})
}

// onConfigMapDeleted is run when an configmap is deleted from the namespace.
func (c *CaddyController) onConfigMapDeleted(obj *v1.ConfigMap) {
	c.syncQueue.Add(ConfigMapDeletedAction{
		resource: obj,
	})
}

func (r ConfigMapAddedAction) handle(c *CaddyController) error {
	logrus.Info("New configmap detected, updating Caddy config...")

	cfg, err := k8s.ParseConfigMap(r.resource)
	if err == nil {
		c.resourceStore.ConfigMap = cfg
	}
	return err
}

func (r ConfigMapUpdatedAction) handle(c *CaddyController) error {
	logrus.Info("ConfigMap resource update detected, updating Caddy config...")

	cfg, err := k8s.ParseConfigMap(r.resource)
	if err == nil {
		c.resourceStore.ConfigMap = cfg
	}
	return err
}

func (r ConfigMapDeletedAction) handle(c *CaddyController) error {
	logrus.Info("ConfigMap resource deletion detected, updating Caddy config...")

	c.resourceStore.ConfigMap = nil
	return nil
}
