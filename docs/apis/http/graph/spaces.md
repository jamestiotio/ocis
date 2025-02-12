---
title: Spaces
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/apis/http/graph
geekdocFilePath: spaces.md
---

{{< toc >}}

## Spaces API

The Spaces API is implementing a subset of the functionality of the
[MS Graph Drives resource](https://learn.microsoft.com/en-us/graph/api/resources/drive?view=graph-rest-1.0).

### Example Space

The JSON representation of a Drive, as handled by the Spaces API, looks like this:
````json
{
  "driveAlias": "project/mars",
  "driveType": "project",
  "id": "storage-users-1$89ad5ad2-5fdb-4877-b8c9-601a9670b925",
  "lastModifiedDateTime": "2023-01-24T21:19:26.417055+01:00",
  "name": "Mars",
  "owner": {
    "user": {
      "displayName": "",
      "id": "89ad5ad2-5fdb-4877-b8c9-601a9670b925"
    }
  },
  "quota": {
    "remaining": 999853685,
    "state": "normal",
    "total": 1000000000,
    "used": 146315
  },
  "root": {
    "eTag": "\"910af0061161c42d8d1224df6c4a2527\"",
    "id": "storage-users-1$89ad5ad2-5fdb-4877-b8c9-601a9670b925",
    "permissions": [
      {
        "grantedToIdentities": [
          {
            "user": {
              "displayName": "Admin",
              "id": "some-admin-user-id-0000-000000000000"
            }
          }
        ],
        "roles": [
          "manager"
        ]
      }
    ],
    "webDavUrl": "https://localhost:9200/dav/spaces/storage-users-1$89ad5ad2-5fdb-4877-b8c9-601a9670b925"
  },
  "special": [
    {
      "eTag": "\"f97829324f63ce778095334cfeb0097b\"",
      "file": {
        "mimeType": "image/jpeg"
      },
      "id": "storage-users-1$89ad5ad2-5fdb-4877-b8c9-601a9670b925!40171bea-3263-47a8-80ef-0ca20c37f45a",
      "lastModifiedDateTime": "2022-02-15T17:11:50.000000496+01:00",
      "name": "Mars_iStock-MR1805_20161221.jpeg",
      "size": 146250,
      "specialFolder": {
        "name": "image"
      },
      "webDavUrl": "https://localhost:9200/dav/spaces/storage-users-1$89ad5ad2-5fdb-4877-b8c9-601a9670b925%2189ad5ad2-5fdb-4877-b8c9-601a9670b925/.space/Mars_iStock-MR1805_20161221.jpeg"
    },
    {
      "eTag": "\"ff38b31d8f109a4fbb98ab34499a3379\"",
      "file": {
        "mimeType": "text/markdown"
      },
      "id": "storage-users-1$89ad5ad2-5fdb-4877-b8c9-601a9670b925!e2167612-7578-46e2-8ed7-971481037bc1",
      "lastModifiedDateTime": "2023-01-24T21:10:23.661841+01:00",
      "name": "readme.md",
      "size": 65,
      "specialFolder": {
        "name": "readme"
      },
      "webDavUrl": "https://localhost:9200/dav/spaces/storage-users-1$89ad5ad2-5fdb-4877-b8c9-601a9670b925%2189ad5ad2-5fdb-4877-b8c9-601a9670b925/.space/readme.md"
    }
  ],
  "webUrl": "https://localhost:9200/f/storage-users-1$89ad5ad2-5fdb-4877-b8c9-601a9670b925"
}
````

## Creating Spaces

### Create a single space `POST /drives`

Create a new space with properties.

{{< tabs "create-space" >}}
{{< tab "Request" >}}

```shell
curl -L -X POST 'https://localhost:9200/graph/v1.0/drives/' \
-H 'Content-Type: application/json' \
--data-raw '{
    "Name": "Marketing",
    "description": "Marketing team resources",
    "quota": {
        "total": 5368709120
    }
}'
```
{{< /tab >}}
{{< tab "Response - 201 Created" >}}

```json {hl_lines=[2,7,15]}
{
  "description":"Marketing team resources",
  "driveAlias":"project/marketing",
  "driveType":"project",
  "id":"storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff",
  "lastModifiedDateTime":"2023-01-18T17:13:48.385204589+01:00",
  "name":"Marketing",
  "owner":{
    "user":{
      "displayName":"",
      "id":"535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
    }
  },
  "quota":{
    "total":5368709120
  },
  "root":{
    "eTag":"\"f91e56554fd9305db81a93778c0fae96\"",
    "id":"storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff",
    "permissions":[
      {
        "grantedToIdentities":[
          {
            "user":{
              "displayName":"Admin",
              "id":"some-admin-user-id-0000-000000000000"
            }
          }
        ],
        "roles":[
          "manager"
        ]
      }
    ],
    "webDavUrl":"https://localhost:9200/dav/spaces/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
  },
  "webUrl":"https://localhost:9200/f/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
}
```
{{< /tab >}}
{{< /tabs >}}

## Reading Spaces

```shell
GET https://ocis.url/graph/{version}/{me/}drives/?{query-parameters}
```

| Component          | Description                                                                                                            |
|--------------------|------------------------------------------------------------------------------------------------------------------------|
| {version}          | The version of the LibreGraph API used by the client.                                                                  |
| {/me}              | The `me` component of the part is optional. If used, you only see spaces where the acting user is a regular member of. |
| {query-parameters} | Optional parameters for the request to customize the response.                                                         |

### List all spaces `GET /drives`

Returns a list of all available spaces, even ones where the acting user is not a regular member of. You need elevated permissions to do list all spaces. If you don't have the elevated permissions, the result is the same like `GET /me/drives`.


{{< hint type=info title="Multiple Administration Personas" >}}

The ownCloud spaces concept draws a strict line between users which can work with the content of a space and others who have the permission to manage the space. A user which is able to manage quota and space metadata does not necessarily need to be able to access the content of a space.

**Space Admin**\
There is a global user role "Space Admin" which grants users some global permissions to manage space quota and some space metadata. This Role enables the user also to disable, restore and delete spaces. He cannot manage space members.

**Space Manager**\
The "Space Manager" is a user which is a regular member of a space because he has been invited. In addition to being part of a space the user can also manage the memberships of the space.

{{< /hint >}}

### List My Spaces `GET /me/drives`

Returns a list of all spaces where the user is a regular member of.

{{< tabs "list-drives" >}}
{{< tab "Request" >}}
`curl -L -k -X GET 'https://localhost:9200/graph/v1.0/me/drives/'`
{{< /tab >}}
{{< tab "Response - 200 OK" >}}
```json
{
    "value": [
        {
            "driveAlias": "personal/admin",
            "driveType": "personal",
            "id": "storage-users-1$some-admin-user-id-0000-000000000000",
            "lastModifiedDateTime": "2023-01-17T21:37:17.692033183+01:00",
            "name": "Admin",
            "owner": {
                "user": {
                    "displayName": "",
                    "id": "some-admin-user-id-0000-000000000000"
                }
            },
            "quota": {
                "remaining": 4497686528,
                "state": "normal",
                "total": 0,
                "used": 0
            },
            "root": {
                "eTag": "\"4b65d2cbce79b2ecc45f876cc6e460d8\"",
                "id": "storage-users-1$some-admin-user-id-0000-000000000000",
                "webDavUrl": "https://localhost:9200/dav/spaces/storage-users-1$some-admin-user-id-0000-000000000000"
            },
            "webUrl": "https://localhost:9200/f/storage-users-1$some-admin-user-id-0000-000000000000"
        },
        {
            "driveAlias": "virtual/shares",
            "driveType": "virtual",
            "id": "a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668",
            "name": "Shares",
            "quota": {
                "remaining": 0,
                "state": "exceeded",
                "total": 0,
                "used": 0
            },
            "root": {
                "eTag": "DECAFC00FEE",
                "id": "a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668",
                "webDavUrl": "https://localhost:9200/dav/spaces/a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668"
            },
            "webUrl": "https://localhost:9200/f/a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668"
        }
    ]
}
```
{{< /tab >}}
{{< /tabs >}}

### List my spaces with filters `GET /me/drives/$filter=<attribute> <comparison operator> <value>`

Returns a list of all project spaces where the user is a regular member of. The filter query parameter supports the `eq`(equals) comparison operator. Possible filter attributes are:

| Attribute | Description                                                                    |
|-----------|--------------------------------------------------------------------------------|
| driveType | The space type. Values could be `project`, `personal`, `virtual`, `mountpoint` |
| id        | The space id. The value needs to be a space ID                                 |

#### Example filters

{{< tabs "filter-drives" >}}
{{< tab "Filter by space type project" >}}

`curl -L -k -X GET 'https://localhost:9200/graph/v1.0/me/drives/$filter=driveType eq project'`

This returns a list of spaces of type `project`.

{{< /tab >}}
{{< tab "Filter by space type mountpoint" >}}

`curl -L -k -X GET 'https://localhost:9200/graph/v1.0/me/drives/$filter=driveType eq mountpoint'`

This returns a list of spaces of type `mountpoint`.

{{< /tab >}}
{{< tab "Filter by space id" >}}

`curl -L -k -X GET 'https://localhost:9200/graph/v1.0/me/drives/$filter=id eq some-space-id'`

This returns one space with the id from the filter if it exists.
{{< /tab >}}
{{< /tabs >}}

### List my spaces with ordering `GET /me/drives/$orderby=<key> <order>`

Returns a list of all spaces ordered by the specified key. Possible order keys are

| key                   | Description                                  |
|-----------------------|----------------------------------------------|
| name	                 | The space name                               |
| lastModifiedDateTime	 | The last modified date and time of the space |

Both sort orders are supported, `asc` (ascending) and `desc` (descending).

#### Ordering examples

{{< tabs "sort-drives" >}}
{{< tab "Sort by name" >}}

`curl -L -k -X GET 'https://localhost:9200/graph/v1.0/me/drives/$orderby=name asc'`

This returns a list of spaces ordered by name ascending.

{{< /tab >}}
{{< tab "Sort by date" >}}

`curl -L -k -X GET 'https://localhost:9200/graph/v1.0/me/drives/$orderby=lastModifiedDateTime desc'`

This returns a list of spaces ordered by the last modified date starting with the latest one.
{{< /tab >}}
{{< /tabs >}}

## Modifying Spaces

Modify the properties of a space. You need elevated permissions to execute this request.

### Set the space quota to 5GB `PATCH /drives/{drive-id}`

To limit the quota of a space you need to set the `quota.total` value. The API response will give back all actual quota properties.

````json
{
  "quota": {
    "remaining": 5368709120,
    "state": "normal",
    "total": 5368709120,
    "used": 0
  }
}
````

| Attribute | Description                                                                                                                                                                                               |
|-----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| remaining | The remaining disk space in `bytes`. If the quota is not limited, this will show the total available disk space.                                                                                          |
| state     | The state of the space in regards to quota usage. This can be used for visual indicators. It can be `normal`(<75%), `nearing`(between 75% and 89%), `critical`(between 90% and 99%) and `exceeded`(100%). |
| total     | The space id. The value needs to be a space ID.                                                                                                                                                           |
| used      | The used disk space in bytes.                                                                                                                                                                             |

{{< tabs "set-space-quota" >}}
{{< tab "Request" >}}
```shell
curl -L -k -X PATCH 'https://localhost:9200/graph/v1.0/drives/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff' \
-H 'Content-Type: application/json' \
--data-raw '{
    "quota": {
        "total": 5368709120
    }
}'
```
{{< /tab >}}
{{< tab "Response - 200 OK" >}}
````json {hl_lines=[17]}
{
    "description": "Marketing team resources",
    "driveAlias": "project/marketing",
    "driveType": "project",
    "id": "storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff",
    "lastModifiedDateTime": "2023-01-18T17:13:48.385204589+01:00",
    "name": "Marketing",
    "owner": {
        "user": {
            "displayName": "",
            "id": "535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
        }
    },
    "quota": {
        "remaining": 5368709120,
        "state": "normal",
        "total": 5368709120,
        "used": 0
    },
    "root": {
        "eTag": "\"f91e56554fd9305db81a93778c0fae96\"",
        "id": "storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff",
        "permissions": [
            {
                "grantedToIdentities": [
                    {
                        "user": {
                            "displayName": "Admin",
                            "id": "some-admin-user-id-0000-000000000000"
                        }
                    }
                ],
                "roles": [
                    "manager"
                ]
            }
        ],
        "webDavUrl": "https://localhost:9200/dav/spaces/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
    },
    "webUrl": "https://localhost:9200/f/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
}
````
{{< /tab >}}
{{< /tabs >}}

### Change the space name, subtitle and alias `PATCH /drives/{drive-id}`

You can change multiple space properties in one request as long as you submit a valid JSON body. Please be aware that some properties need different permissions.

{{< tabs "change-space-props" >}}
{{< tab "Request" >}}
```shell
curl -L -k -X PATCH 'https://localhost:9200/graph/v1.0/drives/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff' \
-H 'Content-Type: application/json' \
--data-raw '{
    "name": "Mars",
    "description": "Mission to mars",
    "driveAlias": "project/mission-to-mars"
}'
```
{{< /tab >}}

{{< tab "Response - 200 OK" >}}
````json {hl_lines=[2,3,7]}
{
    "description": "Mission to mars",
    "driveAlias": "project/mission-to-mars",
    "driveType": "project",
    "id": "storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff",
    "lastModifiedDateTime": "2023-01-19T14:17:36.094283+01:00",
    "name": "Mars",
    "owner": {
        "user": {
            "displayName": "",
            "id": "535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
        }
    },
    "quota": {
        "remaining": 15,
        "state": "normal",
        "total": 15,
        "used": 0
    },
    "root": {
        "eTag": "\"f5fee4fdfeedd6f98956500779eee15e\"",
        "id": "storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff",
        "permissions": [
            {
                "grantedToIdentities": [
                    {
                        "user": {
                            "displayName": "Admin",
                            "id": "some-admin-user-id-0000-000000000000"
                        }
                    }
                ],
                "roles": [
                    "manager"
                ]
            }
        ],
        "webDavUrl": "https://localhost:9200/dav/spaces/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
    },
    "webUrl": "https://localhost:9200/f/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
}
````
{{< /tab >}}
{{< /tabs >}}

## Disabling / Deleting Spaces

### Disable a space `DELETE /drives/{drive-id}`

This operation will make the space content unavailable for all space members. No data will be deleted.

{{< tabs "disable-space" >}}
{{< tab "Request" >}}
```shell
curl -L -k -X DELETE 'https://localhost:9200/graph/v1.0/drives/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff/'
```
{{< /tab >}}

{{< tab "Response - 204 No Content" >}}

This response has no body value.

A disabled space will appear in listings with a `root.deleted.state=trashed` property. The space description and the space image will not be readable anymore.

```json {hl_lines=[18,19,20]}
{
    "description": "Marketing team resources",
    "driveAlias": "project/marketing",
    "driveType": "project",
    "id": "storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff",
    "lastModifiedDateTime": "2023-01-19T14:17:36.094283+01:00",
    "name": "Marketing",
    "owner": {
        "user": {
            "displayName": "",
            "id": "535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
        }
    },
    "quota": {
        "total": 15
    },
    "root": {
        "deleted": {
            "state": "trashed"
        },
        "eTag": "\"f5fee4fdfeedd6f98956500779eee15e\"",
        "id": "storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff",
        "permissions": [
            {
                "grantedToIdentities": [
                    {
                        "user": {
                            "displayName": "Admin",
                            "id": "some-admin-user-id-0000-000000000000"
                        }
                    }
                ],
                "roles": [
                    "manager"
                ]
            }
        ],
        "webDavUrl": "https://localhost:9200/dav/spaces/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
    },
    "webUrl": "https://localhost:9200/f/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
}
```

{{< /tab >}}
{{< /tabs >}}

### Restore a space `PATCH /drives/{drive-id}`

This operation will make the space content available again to all members. No content will be changed.

To restore a space, the Header `Restore: T` needs to be set.
{{< tabs "restore-space" >}}
{{< tab "Request" >}}

```shell
curl -L -X PATCH 'https://localhost:9200/graph/v1.0/drives/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff/' \
-H 'Restore: T' \
-H 'Content-Type: text/plain' \
--data-raw '{}'
```

{{< hint type=info title="Body value" >}}

This request needs an empty body (--data-raw '{}') to fulfil the standard libregraph specificiation even when the body is not needed.

{{< /hint >}}
{{< /tab >}}

{{< tab "Response - 200 OK" >}}

```json
{
    "description": "Marketing team resources",
    "driveAlias": "project/marketing",
    "driveType": "project",
    "id": "storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff",
    "lastModifiedDateTime": "2023-01-19T14:17:36.094283+01:00",
    "name": "Marketing",
    "owner": {
        "user": {
            "displayName": "",
            "id": "535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
        }
    },
    "quota": {
        "remaining": 15,
        "state": "normal",
        "total": 15,
        "used": 0
    },
    "root": {
        "eTag": "\"f5fee4fdfeedd6f98956500779eee15e\"",
        "id": "storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff",
        "permissions": [
            {
                "grantedToIdentities": [
                    {
                        "user": {
                            "displayName": "Admin",
                            "id": "some-admin-user-id-0000-000000000000"
                        }
                    }
                ],
                "roles": [
                    "manager"
                ]
            }
        ],
        "webDavUrl": "https://localhost:9200/dav/spaces/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
    },
    "webUrl": "https://localhost:9200/f/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff"
}
```

{{< /tab >}}
{{< /tabs >}}

### Permanently delete a space `DELETE /drives/{drive-id}`

This operation will delete a space and all its data permanently. This is restricted to spaces which are already disabled.

To delete a space, the Header `Purge: T` needs to be set.

{{< tabs "delete-space" >}}
{{< tab "Request" >}}

```shell {hl_lines=[2]}
curl -L -X DELETE 'https://localhost:9200/graph/v1.0/drives/storage-users-1$535aa42d-a3c7-4329-9eba-5ef48fcaa3ff' \
-H 'Purge: T'
```

{{< hint type=warning title="Data will be deleted" >}}

This request will delete a space and all its content permanently. This operation cannot be reverted.

{{< /hint >}}

{{< /tab >}}
{{< tab "Response - 204 No Content" >}}

This response has no body value.

{{< /tab >}}
{{< tab "Response - 400 Bad Request" >}}

The space to be deleted was not disabled before.

```json
{
    "error": {
        "code": "invalidRequest",
        "innererror": {
            "date": "2023-01-24T19:57:19Z",
            "request-id": "f62af40f-bc18-475e-acd7-e9008d6bd326"
        },
        "message": "error: bad request: can't purge enabled space"
    }
}
```
{{< /tab >}}
{{< /tabs >}}
