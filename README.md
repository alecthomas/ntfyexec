# ntfyexec - send a notification through ntfy.sh if a command fails

```
ntfyexec --title "ZFS pool scrub on storage failed" my-server-topic zpool scrub -w storage
```

If the command fails, it will send a notification to `my-server-topic` with the given title and the output of the
command as the body.
