package rest

test_auth_login_with_id_password_allowed {
    allowed with input as {
        "api": {
            "url": "/api/auth",
            "method": "POST"
        }
    }
}

test_auth_login_with_dooray_allowed {
    allowed with input as {
        "api": {
            "url": "/api/auth/dooray",
            "method": "POST"
        }
    }
}

test_auth_login_with_google_workspace_allowed {
    allowed with input as {
        "api": {
            "url": "/api/auth/google-workspace",
            "method": "GET"
        }
    }
}

test_auth_token_refresh_allowed {
    allowed with input as {
        "api": {
            "url": "/api/auth/token/refresh",
            "method": "POST"
        }
    }
}

test_access_control_permissions_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.read"]
        },
        "api": {
            "url": "/api/access-control/permissions",
            "method": "GET"
        }
    }
}

test_access_control_permissions_read_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.create", "access-control-permission.update"]
        },
        "api": {
            "url": "/api/access-control/permissions",
            "method": "GET"
        }
    }
}

test_access_control_permission_create_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.create"]
        },
        "api": {
            "url": "/api/access-control/permissions",
            "method": "POST"
        }
    }
}

test_access_control_permission_create_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.read"]
        },
        "api": {
            "url": "/api/access-control/permissions",
            "method": "POST"
        }
    }
}

test_access_control_permission_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.read"]
        },
        "api": {
            "url": "/api/access-control/permissions/:permissionId",
            "method": "GET"
        }
    }
}

test_access_control_permission_read_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.update"]
        },
        "api": {
            "url": "/api/access-control/permissions/:permissionId",
            "method": "GET"
        }
    }
}

test_access_control_permission_update_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.update"]
        },
        "api": {
            "url": "/api/access-control/permissions/:permissionId",
            "method": "PUT"
        }
    }
}

test_access_control_permission_update_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.read"]
        },
        "api": {
            "url": "/api/access-control/permissions/:permissionId",
            "method": "PUT"
        }
    }
}

test_access_control_permission_delete_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.delete"]
        },
        "api": {
            "url": "/api/access-control/permissions/:permissionId",
            "method": "DELETE"
        }
    }
}

test_access_control_permission_delete_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.read"]
        },
        "api": {
            "url": "/api/access-control/permissions/:permissionId",
            "method": "DELETE"
        }
    }
}

test_access_control_permissions_all_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.all"]
        },
        "api": {
            "url": "/api/access-control/permissions",
            "method": "GET"
        }
    }
}

test_access_control_permission_all_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-permission.all"]
        },
        "api": {
            "url": "/api/access-control/permissions/:permissionId",
            "method": "PUT"
        }
    }
}

test_access_control_roles_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.read"]
        },
        "api": {
            "url": "/api/access-control/roles",
            "method": "GET"
        }
    }
}

test_access_control_roles_read_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.update"]
        },
        "api": {
            "url": "/api/access-control/roles",
            "method": "GET"
        }
    }
}

test_access_control_role_create_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.create"]
        },
        "api": {
            "url": "/api/access-control/roles",
            "method": "POST"
        }
    }
}

test_access_control_role_create_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.update"]
        },
        "api": {
            "url": "/api/access-control/roles",
            "method": "POST"
        }
    }
}

test_access_control_role_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.read"]
        },
        "api": {
            "url": "/api/access-control/roles/:roleId",
            "method": "GET"
        }
    }
}

test_access_control_role_read_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.create"]
        },
        "api": {
            "url": "/api/access-control/roles/:roleId",
            "method": "GET"
        }
    }
}

test_access_control_role_update_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.update"]
        },
        "api": {
            "url": "/api/access-control/roles/:roleId",
            "method": "PUT"
        }
    }
}

test_access_control_role_update_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.read"]
        },
        "api": {
            "url": "/api/access-control/roles/:roleId",
            "method": "PUT"
        }
    }
}

test_access_control_role_delete_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.delete"]
        },
        "api": {
            "url": "/api/access-control/roles/:roleId",
            "method": "DELETE"
        }
    }
}

test_access_control_role_delete_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.read"]
        },
        "api": {
            "url": "/api/access-control/roles/:roleId",
            "method": "DELETE"
        }
    }
}

test_access_control_roles_all_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.all"]
        },
        "api": {
            "url": "/api/access-control/roles",
            "method": "GET"
        }
    }
}

test_access_control_role_all_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["access-control-role.all"]
        },
        "api": {
            "url": "/api/access-control/roles/:roleId",
            "method": "DELETE"
        }
    }
}

test_member_siginup_allowed {
    allowed with input as {
        "api": {
            "url": "/api/members",
            "method": "POST"
        }
    }
}

test_members_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": []
        },
        "api": {
            "url": "/api/members",
            "method": "GET"
        }
    }
}

test_member_my_information_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": []
        },
        "api": {
            "url": "/api/members/my",
            "method": "GET"
        }
    }
}

test_member_my_information_read_not_allowed {
    not allowed with input as {
        "api": {
            "url": "/api/members/my",
            "method": "GET"
        }
    }
}

test_member_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["member.read"]
        },
        "api": {
            "url": "/api/members/:id",
            "method": "GET"
        }
    }
}

test_member_read_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["member.update"]
        },
        "api": {
            "url": "/api/members/:id",
            "method": "GET"
        }
    }
}

test_member_assgin_roles_update_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["member.update"]
        },
        "api": {
            "url": "/api/members/:id/assign-roles",
            "method": "PUT"
        }
    }
}

test_member_assgin_roles_update_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["member.read"]
        },
        "api": {
            "url": "/api/members/:id/assign-roles",
            "method": "PUT"
        }
    }
}

test_member_approved_update_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["member.update"]
        },
        "api": {
            "url": "/api/members/:id/approved",
            "method": "PUT"
        }
    }
}

test_member_approved_update_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["member.read"]
        },
        "api": {
            "url": "/api/members/:id/approved",
            "method": "PUT"
        }
    }
}

test_member_rejected_update_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["member.update"]
        },
        "api": {
            "url": "/api/members/:id/rejected",
            "method": "PUT"
        }
    }
}

test_member_rejected_update_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["member.read"]
        },
        "api": {
            "url": "/api/members/:id/rejected",
            "method": "PUT"
        }
    }
}

test_members_search_filters_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": []
        },
        "api": {
            "url": "/api/members/search-filters",
            "method": "GET"
        }
    }
}

test_members_search_filters_read_not_allowed {
    not allowed with input as {
        "api": {
            "url": "/api/members/search-filters",
            "method": "GET"
        }
    }
}

test_organization_create_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.create"]
        },
        "api": {
            "url": "/api/organizations",
            "method": "POST"
        }
    }
}

test_organization_create_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": []
        },
        "api": {
            "url": "/api/organizations",
            "method": "POST"
        }
    }
}

test_organizations_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.read"]
        },
        "api": {
            "url": "/api/organizations",
            "method": "GET"
        }
    }
}

test_organization_create_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.create"]
        },
        "api": {
            "url": "/api/organizations",
            "method": "GET"
        }
    }
}

test_organization_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.read"]
        },
        "api": {
            "url": "/api/organizations/:organizationId",
            "method": "GET"
        }
    }
}

test_organization_read_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.create"]
        },
        "api": {
            "url": "/api/organizations/:organizationId",
            "method": "GET"
        }
    }
}

test_organization_delete_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.delete"]
        },
        "api": {
            "url": "/api/organizations/:organizationId",
            "method": "DELETE"
        }
    }
}

test_organization_delete_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.read"]
        },
        "api": {
            "url": "/api/organizations/:organizationId",
            "method": "DELETE"
        }
    }
}

test_organization_change_name_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.update"]
        },
        "api": {
            "url": "/api/organizations/:organizationId/name",
            "method": "PUT"
        }
    }
}

test_organization_change_name_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.read"]
        },
        "api": {
            "url": "/api/organizations/:organizationId/name",
            "method": "PUT"
        }
    }
}

test_organization_change_position_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.update"]
        },
        "api": {
            "url": "/api/organizations/:organizationId/change-position",
            "method": "PUT"
        }
    }
}

test_organization_change_position_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.read"]
        },
        "api": {
            "url": "/api/organizations/:organizationId/change-position",
            "method": "PUT"
        }
    }
}

test_organization_assign_roles_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.update"]
        },
        "api": {
            "url": "/api/organizations/:organizationId/assign-roles",
            "method": "PUT"
        }
    }
}

test_organization_assign_roles_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.read"]
        },
        "api": {
            "url": "/api/organizations/:organizationId/assign-roles",
            "method": "PUT"
        }
    }
}

test_organization_assign_members_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.update"]
        },
        "api": {
            "url": "/api/organizations/:organizationId/assign-members",
            "method": "PUT"
        }
    }
}

test_organization_assign_members_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["organization.read"]
        },
        "api": {
            "url": "/api/organizations/:organizationId/assign-members",
            "method": "PUT"
        }
    }
}

test_site_settings_read_allowed {
    allowed with input as {
        "api": {
            "url": "/api/site/settings",
            "method": "GET"
        }
    }
}

test_site_settings_dooray_login_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["site-settings.read"]
        },
        "api": {
            "url": "/api/site/settings/dooray-login",
            "method": "GET"
        }
    }
}

test_site_settings_dooray_login_read_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": []
        },
        "api": {
            "url": "/api/site/settings/dooray-login",
            "method": "GET"
        }
    }
}

test_site_settings_dooray_login_update_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["site-settings.update"]
        },
        "api": {
            "url": "/api/site/settings/dooray-login",
            "method": "PUT"
        }
    }
}

test_site_settings_dooray_login_update_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["site-settings.read"]
        },
        "api": {
            "url": "/api/site/settings/dooray-login",
            "method": "PUT"
        }
    }
}

test_site_settings_google_workspace_login_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["site-settings.read"]
        },
        "api": {
            "url": "/api/site/settings/google-workspace-login",
            "method": "GET"
        }
    }
}

test_site_settings_google_workspace_login_read_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["site-settings.update"]
        },
        "api": {
            "url": "/api/site/settings/google-workspace-login",
            "method": "GET"
        }
    }
}

test_site_settings_google_workspace_login_update_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["site-settings.update"]
        },
        "api": {
            "url": "/api/site/settings/google-workspace-login",
            "method": "PUT"
        }
    }
}

test_site_settings_google_workspace_login_update_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["site-settings.read"]
        },
        "api": {
            "url": "/api/site/settings/google-workspace-login",
            "method": "PUT"
        }
    }
}

test_site_settings_app_version_read_allowed {
    allowed with input as {
        "api": {
            "url": "/api/site/settings/app-version",
            "method": "GET"
        }
    }
}

test_site_settings_app_version_update_allowed {
    allowed with input as {
        "api": {
            "url": "/api/site/settings/app-version",
            "method": "PUT"
        }
    }
}

test_web_hook_create_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["web-hook.create"]
        },
        "api": {
            "url": "/api/web-hooks",
            "method": "POST"
        }
    }
}

test_web_hook_create_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["web-hook.read"]
        },
        "api": {
            "url": "/api/web-hooks",
            "method": "POST"
        }
    }
}

test_web_hooks_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["web-hook.read"]
        },
        "api": {
            "url": "/api/web-hooks",
            "method": "GET"
        }
    }
}

test_web_hooks_read_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": []
        },
        "api": {
            "url": "/api/web-hooks",
            "method": "GET"
        }
    }
}

test_web_hook_read_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": []
        },
        "api": {
            "url": "/api/web-hooks/:id",
            "method": "GET"
        }
    }
}

test_web_hook_read_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["web-hook.read"]
        },
        "api": {
            "url": "/api/web-hooks/:id",
            "method": "GET"
        }
    }
}

test_web_hook_delete_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["web-hook.delete"]
        },
        "api": {
            "url": "/api/web-hooks/:id",
            "method": "DELETE"
        }
    }
}

test_web_hook_delete_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["web-hook.read"]
        },
        "api": {
            "url": "/api/web-hooks/:id",
            "method": "DELETE"
        }
    }
}

test_web_hook_update_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["web-hook.update"]
        },
        "api": {
            "url": "/api/web-hooks/:id",
            "method": "PUT"
        }
    }
}

test_web_hook_update_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["web-hook.read"]
        },
        "api": {
            "url": "/api/web-hooks/:id",
            "method": "PUT"
        }
    }
}

test_web_hook_note_create_allowed {
    allowed with input as {
        "member": {
            "id": 1,
            "permissions": ["web-hook-note.create"]
        },
        "api": {
            "url": "/api/web-hooks/:id/note",
            "method": "POST"
        }
    }
}

test_web_hook_note_create_not_allowed {
    not allowed with input as {
        "member": {
            "id": 1,
            "permissions": []
        },
        "api": {
            "url": "/api/web-hooks/:id/note",
            "method": "POST"
        }
    }
}