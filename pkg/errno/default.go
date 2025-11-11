package errno

// 一些常用的的错误, 如果你懒得单独定义也可以直接使用
var (
	Success = NewErrNo(SuccessCode, "成功")

	ParamVerifyError  = NewErrNo(ParamVerifyErrorCode, "参数验证失败")
	ParamMissingError = NewErrNo(ParamMissingErrorCode, "缺少必要参数")

	AuthInvalid             = NewErrNo(AuthInvalidCode, "身份验证失败")
	AuthAccessExpired       = NewErrNo(AuthAccessExpiredCode, "令牌已过期")
	AuthNoToken             = NewErrNo(AuthNoTokenCode, "缺少令牌")
	AuthNoOperatePermission = NewErrNo(AuthNoOperatePermissionCode, "没有操作权限")

	InternalServiceError = NewErrNo(InternalServiceErrorCode, "内部服务错误")
	OSOperationError     = NewErrNo(OSOperateErrorCode, "操作系统调用失败")
	IOOperationError     = NewErrNo(IOOperateErrorCode, "输入输出操作失败")

	QiNiuYunFileError = NewErrNo(QiNiuYunFileErrorCode, "七牛云操作失败")

	//  User Module Errors
	UserPasswordIncorrectError = NewErrNo(UserPasswordIncorrect, "密码不正确")
)
