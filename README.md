# terrarium-bot-v2
Manage a Terrarium. Temperature, humidity, heating, misting, lighting all fully configurable and customisable via YAML.  
  
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

## Configuration

[The configuration](cmd/terrarium-bot-v2/configuration.yaml) should be fully customisable and hopefully easy to understand, everything is trigger/action based e.g.

I have a temperature sensor configured:-
```yaml
sensor:
  - id: temperature
    url: http://terrarium-temp-sensor
    jsonPath: Temperature
```
A heating switch configured:-
```yaml
switch:
  - id: heater
    on: http://meross-lan-api/turnOn/heat
    off: http://meross-lan-api/turnOff/heat
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

## The End
Hopefully the configuration should be pretty self explanitory, if you get stuck or there are any features you think might be missing then feel free to create an issue :slightly_smiling_face:.
