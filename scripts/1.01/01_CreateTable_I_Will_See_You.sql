CREATE TABLE Auth.Permissions (
                                  Id UUID PRIMARY KEY,
                                  TenantId UUID NOT NULL,
                                  Name TEXT NOT NULL,
                                  Description TEXT NOT NULL,
                                  TimestampCreated TIMESTAMP without time zone default (now() at time zone 'utc'),
                                  TimestampArchived TIMESTAMP without time zone default NULL
);

CREATE TABLE Auth.Groups (
                             Id UUID PRIMARY KEY,
                             TenantId UUID NOT NULL,
                             Name TEXT NOT NULL,
                             Description TEXT NOT NULL,
                             TimestampCreated TIMESTAMP without time zone default (now() at time zone 'utc')
);

CREATE TABLE Auth.GroupPermissions (
                                       GroupId UUID NOT NULL,
                                       PermissionId UUID NOT NULL,
                                       TimestampCreated TIMESTAMP without time zone default (now() at time zone 'utc')
);

CREATE TABLE Auth.GroupMembers (
                                   GroupId UUID NOT NULL,
                                   UserId UUID NOT NULL,
                                   TimestampCreated TIMESTAMP without time zone default (now() at time zone 'utc')
);

CREATE TABLE Auth.Tenants (
                              Id UUID PRIMARY KEY,
                              Name TEXT NOT NULL,
                              TimestampCreated TIMESTAMP without time zone default (now() at time zone 'utc')
);

CREATE TABLE Auth.TenantUsers (
                                  TenantId UUID NOT NULL,
                                  UserId UUID NOT NULL,
                                  TimestampCreated TIMESTAMP without time zone default (now() at time zone 'utc')
);
