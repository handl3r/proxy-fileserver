package middlewares

type MiddlewareProvider interface {
	GetAuthorizationProcessor() *AuthorizationProcessor
}

type middlewareProviderImpl struct {
	authorizationProcessor *AuthorizationProcessor
}

func NewMiddlewareProvider(publicKeyLocation string) MiddlewareProvider {
	return middlewareProviderImpl{
		authorizationProcessor: NewAuthorizationProcessor(publicKeyLocation),
	}
}

func (p middlewareProviderImpl) GetAuthorizationProcessor() *AuthorizationProcessor {
	return p.authorizationProcessor
}
