package middlewares

type MiddlewareProvider interface {
	GetAuthorizationProcessor() *AuthorizationProcessor
}

type middlewareProviderImpl struct {
	authorizationProcessor *AuthorizationProcessor
}

func NewMiddlewareProvider(publicKey string) MiddlewareProvider {
	return middlewareProviderImpl{
		authorizationProcessor: NewAuthorizationProcessor(publicKey),
	}
}

func (p middlewareProviderImpl) GetAuthorizationProcessor() *AuthorizationProcessor {
	return p.authorizationProcessor
}
