# CSV Run Extractor

CLI tool to extract full throttle runs from automotive logs. The app works by taking a configured throttle parameter from the CSV and identifies a series of consecutive full throttle entries - these are then grouped and extracted into separate files for each full throttle run.

In addition, the tool will also filter out configured fields and order them accordingly, negating the need to manually select fields that are relevant to the run


## Usage
 - Download the executable and config file from the [release page](https://github.com/ftj2006/csv-run-extractor/releases).
 - Simply drag and drop your CSV file onto the executable and it will generate the extracted runs in files with the same name followed by "Run X"

## Configuration
The [config.yaml](/config.yaml) is used to determine what field is used to measure full throttle, how many occurences are required to meet "Run" criteria and what minimum percentage they should be

```yaml
runLimits: 
  throttleField: "Acceleration pedal position (%)" #The field used to determine full throttle status
  minThrottleValue: 95 #Minimum value needed in the throttle field
  minCount: 20 #Amount of consecutive lines needed to count as a run
fields: #Full list of fields to output in the generated files. The order used here will determine the order in the generated files.
  - "Time (ISO)"
  - "..."
```
