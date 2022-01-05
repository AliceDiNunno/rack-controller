package kubernetes

import (
	"context"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8s kubernetesInstance) GetSecretsList(namespace string) ([]string, *e.Error) {
	secrets, err := k8s.Client.CoreV1().Secrets(namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, e.Wrap(err).Append(ErrUnableToGetRessource)
	}

	var secretNames []string
	for _, secret := range secrets.Items {
		secretNames = append(secretNames, secret.Name)
	}

	return secretNames, nil
}

func (k8s kubernetesInstance) GetSecret(namespace string, name string) (domain.Secret, *e.Error) {
	data, err := k8s.Client.CoreV1().Secrets(namespace).Get(context.Background(), name, v1.GetOptions{})

	if err != nil {
		return domain.Secret{}, e.Wrap(err).Append(ErrSecretNotFound)
	}

	var env = []domain.Environment{}

	for envKey, envValue := range data.Data {
		env = append(env, domain.Environment{
			Name:  envKey,
			Value: string(envValue),
		})
	}

	return domain.Secret{
		Name:    data.Name,
		Content: env,
	}, nil
}

func (k8s kubernetesInstance) CreateSecret(namespace string, request domain.SecretCreationRequest) *e.Error {
	secret := &v12.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      request.Name,
			Namespace: namespace,
		},
	}

	_, err := k8s.Client.CoreV1().Secrets(namespace).Create(context.Background(), secret, v1.CreateOptions{})
	if err != nil {
		return e.Wrap(err).Append(ErrSecretCreationFailed)
	}

	return nil
}

func (k8s kubernetesInstance) DeleteSecret(namespace string, name string) *e.Error {
	err := k8s.Client.CoreV1().Secrets(namespace).Delete(context.Background(), name, v1.DeleteOptions{})
	if err != nil {
		return e.Wrap(err).Append(ErrSecretDeletionFailed)
	}

	return nil
}

func (k8s kubernetesInstance) UpdateSecret(namespace string, name string, request domain.SecretUpdateRequest) *e.Error {
	secret, err := k8s.Client.CoreV1().Secrets(namespace).Get(context.Background(), name, v1.GetOptions{})
	if err != nil {
		return e.Wrap(err).Append(ErrSecretNotFound)
	}

	secret.Data = make(map[string][]byte)

	for _, env := range request.Content {
		secret.Data[env.Name] = []byte(env.Value)
	}

	_, err = k8s.Client.CoreV1().Secrets(namespace).Update(context.Background(), secret, v1.UpdateOptions{})
	if err != nil {
		return e.Wrap(err).Append(ErrSecretUpdateFailed)
	}

	return nil
}
