package rest

import "golang.org/x/net/context"

const pathParams = "PATH_PARAMS"

func PathParamFromContext(ctx context.Context) map[string]string {
	return *(ctx.Value(pathParams).(*map[string]string))
}
