wp-cookie-tail
==============
Configure your Wordpress Apache server like this to get a response log (do not use the request log...)

```
<IfModule log_config_module>
    LogFormat "%h %l %u %t \"%r\" %>s %b <strong>\"%{set-cookie}o\"</strong>" common2
    CustomLog "logs/private/wp_cookie_log" common2
</IfModule>
```

These logs can be used with wp-cookie-tail to implement a very simple single sign on system when
used in conjunction with [wp-cookie-verify](https://github.com/acjohnson/wp-cookie-verify).
