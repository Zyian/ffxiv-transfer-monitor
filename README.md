# FFXIV Transfer Monitor

This short Go program scraps [https://na.finalfantasyxiv.com/lodestone/worldstatus/](https://na.finalfantasyxiv.com/lodestone/worldstatus/) every minute to determine the status of the servers.

Once it finds that Faerie is open then it spams me on Discord 5 times. Then goes back to monitoring

## Running Code
Why? This is honestly cause I don't want to watch the site overnight.

```bash
$ chmod +x run.sh
$ ./run.sh
```