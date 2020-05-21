# actions-usage
A command-line tool that displays the current GitHub Actions usage of an entire organization. 

### Usage

```
action-usage [organization]
```

### Example

```
./actions-usage quarkusio
Finding workflows running on all repositories on quarkusio......................
Analyzing jobs......

           wf id  queue/  run / comp                  name          event  repo             created    source
      -----------------------------------------------------------------------------------------------------------------------------------------
  1.   111023796    0q /   1r /   0c  Quarkus CI - JDK 8 N       schedule  quarkus          00 h 18 m  quarkusio/quarkus:master
  2.   110986713    0q /   9r /  17c            Quarkus CI   pull_request  quarkus          01 h 20 m  stuartwdouglas/quarkus:webjars-locator
  3.   111007231    7q /  13r /   6c            Quarkus CI   pull_request  quarkus          00 h 48 m  stuartwdouglas/quarkus:amqp-test2
  4.   111001638    3q /  13r /  10c            Quarkus CI   pull_request  quarkus          01 h 00 m  stuartwdouglas/quarkus:7188
  5.   110568528   10q /  14r /   2c            Quarkus CI   pull_request  quarkus          10 h 22 m  gastaldi/quarkus:messagebodywriter
  6.   110947429    0q /   1r /   0c  Native Test - develo       schedule  quarkus-quickst  02 h 13 m  quarkusio/quarkus-quickstarts:master
          Total:   20q /  51r /  35c

```

### Setup

1. [Download](https://github.com/n1hility/actions-usage/releases/tag/v0.1.0), extract and install action-usage from the respective platform archive.
2. Create a GitHub Personal Access Token: https://github.com/settings/tokens
3. You only need a read-only token so deselect all check boxes
4. Cut and paste the token into a token file in your home directory:

```
echo [PASTE HERE] > ~/.actions-usage.tok
```

### Building 

If you choose to build yourself:

1. Install the latest golang: https://golang.org/dl/
2. git clone git@github.com:n1hility/actions-usage
3. go get
4. go build
5. ./actions-usage

