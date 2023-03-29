# Terrarium Bot
[![build](https://github.com/alec-pinson/terrarium-bot/actions/workflows/build.yaml/badge.svg?event=release)](https://github.com/alec-pinson/terrarium-bot/releases/latest)
[![Latest release](https://img.shields.io/github/v/release/alec-pinson/terrarium-bot?label=Latest%20release)](https://github.com/alec-pinson/terrarium-bot/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/alec-pinson/terrarium-bot)](https://github.com/alec-pinson/terrarium-bot/releases/latest)  

**CURRENTLY IN TESTING**  - I have been running this for a week or so now and I am still fixing minor issues.  
  
Manage a Terrarium. Temperature, humidity, heating, misting, lighting, all fully configurable and customisable via YAML.  
  
I currently run this in a Kubernetes cluster (3 Raspberry Pis using [K3S](https://k3s.io/)). I run it along side my other apps:
- [sonoff-lan-api](https://github.com/alec-pinson/sonoff-lan-api) - for controlling sonoff devices over lan
- [meross-lan-api](https://github.com/alec-pinson/meross-lan-api) - for controlling meross devices over lan
- [gpio-to-api](https://github.com/alec-pinson/gpio-to-api) - for controlling gpio and pulling sensor information

## Features
- Fully customisable
- Turn on/off switches
- Pull data from sensors
- Use the data from sensors to trigger an action e.g. switch off/on 
- Send custom alert notifications e.g. humidity low, temp high etc
- Easy to add custom API endpoints that can be called which then run actions (see example below)
- Time based switches, set what should be on during the day, what should be on during the night
- Sunset/sunrise mode, if you want to turn on some lights a few minutes before others in the morning
- Will send a notification if a sensor stops responding
- Will send a notification if a switch stops responding
- Dry run mode using set `DRY_RUN=true` no switches will be turned on or off + no alert notifications will be sent
- Debug mode for extra output - `DEBUG=true`
- Health probe endpoints for Kubernetes `/health/live` and `/health/ready`

## Getting Started
1. Create a [configuration.yaml](#configuration)
2. Run the docker image mounting your config file
```
docker run -v $PWD:/config -e CONFIG_FILE=/config/configuration.yaml ghcr.io/alec-pinson/terrarium-bot:0.0.3
```

## Configuration

[The configuration](cmd/terrarium-bot/configuration.yaml) should be fully customisable and hopefully easy to understand, everything is trigger/action based e.g.

I have a temperature sensor configured:-
```yaml
sensor:
  - id: temperature
    url: http://terrarium-temp-sensor:8080
    jsonPath: Temperature
    unit: °C
```
A heating switch configured:-
```yaml
switch:
  - id: heater
    on: http://meross-lan-api/turnOn/heat-lamp
    off: http://meross-lan-api/turnOff/heat-lamp
    status: http://meross-lan-api:8080/status/heat-lamp
    jsonPath: Status
```
This heater is controlled by the below trigger, e.g. heating on when below 28 in the day, on when below 22 at night otherwise the heating is off:-
```yaml
trigger:
  - id: heating
    sensor: temperature
    when:
      day:
        below: 28
      night:
        below: 22
    action:
      - switch.heater.on
    else:
      - switch.heater.off
```

Multiple actions can be given to a trigger, e.g. turning lights off for 5 minutes before misting (I like to give my geckos a warning they're about to get sprayed in the face :slightly_smiling_face:):-
```yaml
  - id: mist
    sensor: humidity
    when:
      day:
        below: 62 # trigger if humidity drops below 62
        every: 4h # trigger every 4h, if not triggered by the above low humidity (just in case they want a drink!)
    action:
      - switch.small_led.off
      - switch.big_led.off
      - switch.uvb.off
      - switch.small_led.disable # disable other wise 'day mode' will turn these back on
      - switch.big_led.disable # disable other wise 'day mode' will turn these back on
      - switch.uvb.disable # disable other wise 'day mode' will turn these back on
      - switch.fan.disable.45m # there's an action that turns on the fan if humidity is high and it's about to be
      - alert.humidity.disable.1h # disable humidity alert for 1 hour
      - echo.Misting will begin in 5m
      - sleep.5m
      - switch.mister.on.7s
      - sleep.2s
      - switch.small_led.enable
      - switch.big_led.enable
      - switch.uvb.enable
      - switch.small_led.on
      - switch.big_led.on
      - switch.uvb.on
      - sleep.1h # do not trigger this for at least another hour
```

Call an action from a URL endpoint, I use my gpio-to-api app to monitor a gpio button, when pressed it can trigger a call to a URL (I press my button before openning the terrarium doors as humidity drops and triggers misting otherwise):-
```yaml
  - id: preventMist
    endpoint: /prevent-mist
    action:
      - trigger.mist.disable.40m # disable the above trigger for 40m
```

You can also send alerts if too hot or if the humidity is low for example.  
Configure a notification channel (currently only support for pushover, let me know if you want something else):-
```yaml
notification:
  - id: pushover
    device: Alec-Phone
    sound: tugboat
    userToken: PUSHOVER_USER_TOKEN # env variable name
    apiToken: PUSHOVER_APP_TOKEN # env variable name
    antiSpam: 1h # only receive alerts once an hour
```
Configure an alert if temperature gets below/above certain values during the day or night:-
```yaml
alert:
  - id: temperature
    sensor: temperature
    when:
      day:
        below: 22
        above: 32
      night:
        below: 17
        above: 26
    after: 30m # only alert after 30 minutes
    notification:
      - pushover
```

## Example Actions List
| Action                        | Description                                  |
|-------------------------------|----------------------------------------------|
| sleep.5s                      | Pause execution for 5 seconds                |
| echo.Misting will begin in 5m | Write log message 'Misting will begin in 5m' |
| switch.fan.on                 | Turn on the fan                              |
| switch.fan.on.10m             | Turn on the fan for 10 minutes               |
| switch.fan.off                | Turn off the fan                             |
| switch.fan.disable            | Disable the fan control switch               |
| switch.fan.disable.45m        | Disable the fan control switch for 45m       |
| switch.fan.enable             | Enable the fan control switch                |
| trigger.mist.disable          | Disable the mist trigger                     |
| trigger.mist.disable.40m      | Disable the mist trigger for 40 minutes      |
| trigger.mist.enable           | Enable the mist trigger                      |
| alert.humidity.disable        | Disable the humidity alert                   |
| alert.humidity.disable.1h     | Disable the humidity alert for 1 hour        |
| alert.humidity.enable         | Enable the humidity alert                    |

## Example Log Output
```
2023/03/28 14:55:50 Starting: Terrarium bot
2023/03/28 14:55:50 Loading configuration from '/config/configuration.yaml'...
2023/03/28 14:55:50 Configuration loaded...
2023/03/28 14:55:50 Monitoring sensor 'temperature' (28°C)
2023/03/28 14:55:50 Monitoring sensor 'humidity' (72%)
2023/03/28 14:55:55 Starting API server...
2023/03/28 14:55:55 Switch On: 'small_led' (Setting Day Time configuration)
2023/03/28 14:55:56 API Server started...
2023/03/28 14:55:56 Switch On: 'uvb' (Setting Day Time configuration)
2023/03/28 14:55:57 Switch Off: 'fan' (72%/84%)
2023/03/28 14:56:00 Started: Terrarium bot
2023/03/28 14:56:01 Alert: Terrarium bot started
2023/03/28 15:01:57 Switch On: 'heater' (27°C/28°C)
2023/03/28 15:02:57 Switch Off: 'heater' (28°C/28°C)
2023/03/28 18:54:20 Switch On: 'heater' (27°C/28°C)
2023/03/28 18:56:00 Switch Off: 'small_led' (Trigger 'mist' scheduled every 4h0m0s)
2023/03/28 18:56:00 Switch Off: 'big_led' (Trigger 'mist' scheduled every 4h0m0s)
2023/03/28 18:56:00 Switch Off: 'uvb' (Trigger 'mist' scheduled every 4h0m0s)
2023/03/28 18:56:00 Switch Disabled: 'small_led'
2023/03/28 18:56:00 Switch Disabled: 'big_led'
2023/03/28 18:56:00 Switch Disabled: 'uvb'
2023/03/28 18:56:00 Switch Disabled: 'fan' for 45m0s
2023/03/28 18:56:00 Alert Disabled: 'humidity' for 1h0m0s
2023/03/28 18:56:00 Misting will begin in 5m
2023/03/28 19:01:00 Switch On: 'mister' for 7s (Trigger 'mist' scheduled every 4h0m0s)
2023/03/28 19:01:07 Switch Off: 'mister' (7s has elapsed)
2023/03/28 19:01:09 Switch Enabled: 'small_led'
2023/03/28 19:01:09 Switch Enabled: 'big_led'
2023/03/28 19:01:09 Switch Enabled: 'uvb'
2023/03/28 19:01:09 Switch On: 'small_led' (Trigger 'mist' scheduled every 4h0m0s)
2023/03/28 19:01:10 Switch On: 'big_led' (Trigger 'mist' scheduled every 4h0m0s)
2023/03/28 19:01:10 Switch On: 'uvb' (Trigger 'mist' scheduled every 4h0m0s)
2023/03/28 20:44:29 Switch Off: 'heater' (28°C/28°C)
2023/03/28 20:45:29 Switch On: 'heater' (27°C/28°C)
2023/03/28 20:46:29 Switch Off: 'heater' (28°C/28°C)
2023/03/28 20:47:29 Switch On: 'heater' (27°C/28°C)
2023/03/28 20:55:28 Switch Off: 'big_led' (Setting Sunset configuration)
2023/03/28 20:55:28 Alert Disabled: 'temperature' for 1h0m0s
2023/03/28 21:00:29 Switch Off: 'small_led' (Setting Night Time configuration)
2023/03/28 21:00:30 Switch Off: 'uvb' (Setting Night Time configuration)
2023/03/28 21:00:30 Switch On: 'camera-night-vision-1' (Setting Night Time configuration)
2023/03/28 21:00:30 Switch Off: 'heater' (27°C/22°C)
2023/03/28 21:00:30 Switch On: 'camera-night-vision-2' (Setting Night Time configuration)
2023/03/28 21:00:31 Switch On: 'camera-night-vision-3' (Setting Night Time configuration)
2023/03/28 23:00:58 Switch On: 'fan' for 10m (88%/87%)
2023/03/28 23:10:58 Switch Off: 'fan' (10m has elapsed)
2023/03/28 23:21:58 Switch On: 'fan' for 10m (89%/87%)
2023/03/28 23:31:58 Switch Off: 'fan' (10m has elapsed)
2023/03/28 23:42:58 Switch On: 'fan' for 10m (90%/87%)
2023/03/28 23:52:58 Switch Off: 'fan' (10m has elapsed)
2023/03/29 00:03:58 Switch On: 'fan' for 10m (90%/87%)
2023/03/29 00:13:58 Switch Off: 'fan' (10m has elapsed)
2023/03/29 00:24:58 Switch On: 'fan' for 10m (90%/87%)
2023/03/29 00:31:42 Switch On: 'heater' (21°C/22°C)
2023/03/29 00:34:58 Switch Off: 'fan' (10m has elapsed)
2023/03/29 00:40:42 Switch Off: 'heater' (22°C/22°C)
2023/03/29 07:55:12 Trigger Disabled: 'mist' for 35m0s
2023/03/29 07:55:12 Switch On: 'small_led' (Setting sunrise configuration)
2023/03/29 07:55:12 Switch Off: 'camera-night-vision-1' (Setting sunrise configuration)
2023/03/29 07:55:13 Switch Off: 'camera-night-vision-2' (Setting sunrise configuration)
2023/03/29 07:55:14 Switch Off: 'camera-night-vision-3' (Setting sunrise configuration)
2023/03/29 07:57:11 Switch Off: 'heater' (22°C/22°C)
2023/03/29 08:00:11 Switch On: 'heater' (22°C/28°C)
2023/03/29 08:00:15 Switch On: 'uvb' (Setting Day Time configuration)
```

## The End
Hopefully the configuration should be pretty self explanitory, if you get stuck or there are any features you think might be missing then feel free to create an issue :slightly_smiling_face:.
