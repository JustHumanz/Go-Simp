# Roles and Taging

## Check Yourself
```slash
/mytags
```
small{**Role permission required: No permission needed**}
%br%
Show following group or member vtuber in your account (including your following roles).

## Check Roles
```slash
/role-info slash{role-name}
```
small{**Role permission required: No permission needed**}
%br%
Show following group or member vtuber in your roles.

### Example
```slash
/role-info slash{role-name: role{@nijisanji}}
```

## Taging Yourself
```slash
/tag-me slash{vtuber-group} slash{vtuber-name} slash{reminder (per-minutes, optional)}
```
small{**Role permission required: No permission needed**}
%br%
Tag yourself for get any notification schedules and fanarts from one or many group or member vtuber. You can enter slash{vtuber-group} or slash{vtuber-name}, or maybe both, but you can't empty both. You can set slash{reminder} to get notification before started.

### Example 1
```slash
/tag-me slash{vtuber-group: brave}
```
### Example 2
```slash
/tag-me slash{vtuber-name: evelyn}
```
### Example 3
```slash
/tag-me slash{vtuber-group: 774inc} slash{reminder: 20}
```

## Taging Roles
```slash
/tag-role slash{role-name} slash{vtuber-group} slash{vtuber-name} slash{reminder (per-minutes, optional)}
```
small{**Role permission required: Manage Channel or Higher**}