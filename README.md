# Send a notification through ntfy.sh if a command fails

```
ntfyexec --title "ZFS pool scrub on storage failed" my-server-topic zpool scrub -w storage
```

If the command fails, it will send a notification to `my-server-topic` with the given title and the output of the
command as the body.

## Usage

```
Usage: ntfyexec <topic> <command> ... [flags]

Execute a command, and send a notification to ntfy.sh if it fails.

The combined stdout and stderr will be sent as the notification body.

Arguments:
  <topic>          Ntfy topic.
  <command> ...    Command to execute.

Flags:
  -h, --help            Show context-sensitive help.
      --title=STRING    Ntfy notification title, if any ($NTFY_TITLE).
      --token=STRING    Ntfy access token, if any ($NTFY_TOKEN).
```
