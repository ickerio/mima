# mima
**MI**cro **MA**nager for games (and services) on your VPS

mima is a small command line interface tool that sits on your personal computer and allows you to rapidly start and stop game servers and other services when not in use to save money. We found that we'd only ever play on our community minecraft server for several hours of the week, despite this we'd be paying (per hour) for the whole week. Logging onto the VPS provider, configuring a new server, manually transfering a massive 500mb world save, then downloading java and ssh-ing into the server to start it seemed infeasible if we'd want to cut back on these unnecessary 161 hosting hours. mima aims to provide a lightweight, easily configurable and quick solution to this problem. We calculated it would bring our server cost down from $20/month to merely $0.84/month.

> At current release mima only supports Vultr and defaults to a minecraft server service. We intend to expand on this in the future.

## Installation
See [releases](https://github.com/ickerio/mima/releases) to download the most recent version of the mima CLI executable for your system

## Config
By default a `.mima.yml` file will be loaded but the file name can be configured with the global `--config` flag.
```yml
keys:
  vultr: API_KEY_HERE
  digitalocean: API_KEY_HERE

servers:
  - name: minecraft_smp # This will also appear as the instance / droplet name on the VPS provider
    provider: Vultr # Vultr or DigitalOcean
    plan: 201 
    region: 19
    os: 338
```

Your config file should include any API keys from the services you use. We currently support [Vultr](https://my.vultr.com/settings/#settingsapi) and ~~DigitalOcean~~ as VPS providers. Next, list any services you want mima to manage - these do not have to be already running. Note that the `plan`, `region` and `os` fields are unique IDs for the given provider. Please use the corresponding commands (`plans`, `regions` and `os`) to list human readable forms with their respective IDs.

## Commands

### info
**Usage** `mima info <name of server>` 

**Example** `mima info minecraft_smp`

Outputs server information if specified server is online. The `<name of server>` parameter must correspond to a server name in the config file.

### start / stop
**Usage** `mima <start/stop> <name of server>`

**Example** `mima start minecraft_smp` or `mima stop minecraft_smp`

Starts or stops a server with the specified name and plan, region and OS. The `<name of server>` parameter must correspond to a server name in the config file.

### plans / regions / os
**Usage** `mima <plans/regions/os> <vps provider>`

**Example** `mima plans Vultr` or `mima regions DigitalOcean`

List human readable strings of available plans, regions or operating systems with their respective ID to aid with configuring the YAML file