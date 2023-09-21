package rest

default allowed = false

# 인증/인가가 불필요한 Public API에 대한 허용 정책
allowed {
    requried_permission := data.api[input.api.url][input.api.method]
    count(requried_permission) == 0
}

# 인증된 모든 멤버에게 허용 정책
allowed {
    input.member.id > 0 # 인증된 사용자 체크

    required_permissions := data.api[input.api.url][input.api.method]
    some p
    required_permissions[p] == "all-authenticated-members"
}

# 권한이 있는 멤버에게만 허용 정책
allowed {
    member_permissions := input.member.permissions
    required_permissions := data.api[input.api.url][input.api.method]

    satisfied_permissions := {p | permissionmatch(member_permissions[_], required_permissions[p], ".")}

    count(required_permissions) > 0
    count(satisfied_permissions) == count(required_permissions)
}

permissionmatch(permission, req_permission, delim) = true {
    permission == req_permission
} else = result { # else문으로 여러 규칙 바디를 연결하면 첫 번째 바디의 조건이 만족하지 않았을 때 다음 바디의 조건을 체크하도록 규칙을 작성한다.
    permission_details := split(permission, delim)
    result = count(permission_details) == 2

    req_permission_details := split(req_permission, delim)
    result = count(req_permission_details) == 2

    result = permission_details[1] == "all"
    result = permission_details[0] == req_permission_details[0]
}