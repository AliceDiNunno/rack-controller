package kubernetes

import (
	"context"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s kubernetesInstance) GetConfigMapList(namespace string) ([]string, *e.Error) {
	data, err := k8s.Client.CoreV1().ConfigMaps(namespace).List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	var mapList []string

	for _, configMapEntry := range data.Items {
		mapList = append(mapList, configMapEntry.Name)
	}

	return mapList, nil
}

func (k8s kubernetesInstance) GetConfigMap(namespace string, name string) (domain.ConfigMap, *e.Error) {
	data, err := k8s.Client.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, v1.GetOptions{})

	if err != nil {
		return domain.ConfigMap{}, e.Wrap(err).Append(ErrConfigMapNotFound)
	}

	var env = []domain.Environment{}

	for envKey, envValue := range data.Data {
		env = append(env, domain.Environment{
			Name:  envKey,
			Value: envValue,
		})
	}

	return domain.ConfigMap{
		Name:    data.Name,
		Content: env,
	}, nil
}

func (k8s kubernetesInstance) CreateConfigMap(namespace string, request domain.ConfigMapCreationRequest) *e.Error {
	data := v12.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name: request.Name,
		},
	}

	_, err := k8s.Client.CoreV1().ConfigMaps(namespace).Create(context.Background(), &data, v1.CreateOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrConfigMapCreationFailed)
	}

	return nil
}

func (k8s kubernetesInstance) DeleteConfigMap(namespace string, name string) *e.Error {
	err := k8s.Client.CoreV1().ConfigMaps(namespace).Delete(context.Background(), name, v1.DeleteOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrConfigMapDeletionFailed)
	}

	return nil
}

func (k8s kubernetesInstance) UpdateConfigMap(namespace string, name string, request domain.ConfigMapUpdateRequest) *e.Error {
	data, err := k8s.Client.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, v1.GetOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrConfigMapNotFound)
	}

	data.Data = make(map[string]string)

	for _, env := range request.Content {
		data.Data[env.Name] = env.Value
	}

	_, err = k8s.Client.CoreV1().ConfigMaps(namespace).Update(context.Background(), data, v1.UpdateOptions{})

	if err != nil {
		return e.Wrap(err).Append(ErrConfigMapUpdateFailed)
	}

	return nil
}
