pdf-report
==========

[![Actions Status](https://github.com/BugDiver/pdf-report/workflows/build/badge.svg)](https://github.com/BugDiver/pdf-report/actions)


**Sample HTML Report documemt**

<img src="https://github.com/BugDiver/pdf-report/raw/master/assets/sample.png" alt="Create New Project preview" style="width: 600px;"/>

Installation
------------

```
gauge install pdf-report
```

* Installing specific version
```
gauge install pdf-report --version 2.1.0
```

#### Offline installation
* Download the plugin from [Releases](https://github.com/BugDiver/pdf-report/releases)
```
gauge install pdf-report --file pdf-report-0.0.1-linux.x86_64.zip
```

#### Build from Source

##### Requirements
* [Golang](http://golang.org/)

##### Compiling

Compilation
```

go run build/make.go
```

For cross-platform compilation

```
go run build/make.go --all-platforms
```

##### Installing
After compilation

```
go run build/make.go --install
```


#### Creating distributable

Note: Run after compiling

```
go run build/make.go --distro
```

For distributable across platforms: Windows and Linux for both x86 and x86_64

```
go run build/make.go --distro --all-platforms
```

New distribution details need to be updated in the `pdf-report-install.json` file in the [gauge plugin repository](https://github.com/getgauge/gauge-repository) for a new version update.

Configuration
-------------

The HTML report plugin can be configured by the properties set in the
`env/default.properties` file in the project.

The configurable properties are:

**gauge_reports_dir**

-  Specifies the path to the directory where the execution reports will
   be generated.

-  Should be either relative to the project directory or an absolute
   path. By default it is set to `reports` directory in the project

**overwrite_reports**

-  Set to ``true`` if the reports **must be overwritten** on each
   execution maintaining only the latest execution report.

-  If set to `false` then a _**new report**_ will be generated on each execution in the reports directory in a nested time-stamped directory. By sdefault it is set to `true`.


**gauge_pdf_report_page_orientation**

-  Specifies the orientation for the report pages.
-  Defaults to  Portrait.
-  Accepted options are "P" or "Portrait" and "L" or "Landscape"

