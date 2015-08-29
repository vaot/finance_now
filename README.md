# finance_now
Get realtime stock quotes

# How to use it

```
go install https://github.com/vaot/finance_now

finance_now -quotes=<YOUR SYMBOLS HERE SEPARATE BY COMMA>
```

# With a limit and slack integration

Make sure you have ```SLACK_WEBHOOK_URL``` set in your env.
```
go install https://github.com/vaot/finance_now/hooks_api

hooks_api

finance_now -quotes=tsla, -limits=280
```

As soon as the trade price hit that limit specified, you will get notifications on slack.
