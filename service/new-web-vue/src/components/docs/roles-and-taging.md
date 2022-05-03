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
Tag yourself for get any notification schedules and fanarts from group or member vtuber. You can enter slash{vtuber-group} or slash{vtuber-name}, or maybe both, but you can't empty both. You can set slash{reminder} to get notification before started.

### Example 1

```slash
/tag-me slash{vtuber-group: brave}
```

### Example 2

```slash
/tag-me slash{vtuber-name: evelyn}
```

### Example 3 (Set Reminder)

```slash
/tag-me slash{vtuber-group: 774inc} slash{reminder: 20}
```

## Taging Roles

```slash
/tag-role slash{role-name} slash{vtuber-group} slash{vtuber-name} slash{reminder (per-minutes, optional)}
```

small{**Role permission required: Manage Channel or Higher**}
%br%
Tag roles for get notification from group or member vtuber like previous one.

### Example 1

```slash
/tag-roles slash{role-name: role{@maha5}} slash{vtuber-group: maha5}
```

### Example 2

```slash
/tag-roles slash{role-name: role{@Watame}} slash{vtuber-name: sheep}
```

### Example 3 (Set Reminder)

```slash
/tag-roles slash{role-name: role{@Tenshi}} slash{vtuber-name: tenshi} slash{reminder: 30}
```

## Remove tag for yourself

```slash
/del-tag slash{vtuber-group} slash{vtuber-name}
```

small{**Role permission required: No permission needed**}
%br%
Remove tag for group or member vtuber for yourself. You must fill in ones or both like previous stage.

### Example 1

```slash
/del-tag slash{vtuber-group: kamitsubaki}
```

### Example 2

```slash
/del-tag slash{vtuber-name: gura}
```

## Remove tag for roles

```slash
/del-role slash{role-name} slash{vtuber-group} slash{vtuber-name}
```

small{**Role permission required: Manage Channel or Higher**}
%br%
Remove tag for roles like previous one.

### Example 1

```slash
/del-role slash{role-name: role{@voms}} slash{vtuber-group: voms}
```

### Example 2

```slash
/del-role slash{role-name: role{@doog}} slash{vtuber-name: Koronen}
```
