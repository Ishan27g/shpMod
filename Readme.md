
Cli to rewrite `shipyard module blueprints` to automate enabling/disabling of a combination of shipyard [modules](https://shipyard.run/docs/resources/module) and their dependencies. 
It also wraps calls to `shipyard` by switching dir/ to the target directory containing the modules definitions.

1. Sample shipyard modules blueprint.
> Note : all modules are enabled 

```hcl
module "network" { 
  source = "./network"
}
module "nats" {
  depends_on = ["module.network"]
  source     = "./nats"
}
module "database-mysql-1" {
  depends_on = ["module.network"]
  source     = "./mysql"
}
module "database-mysql-2" {
  depends_on = ["module.network"]
  source     = "./mysql"
}
module "service-1" {
  depends_on = ["module.network", "module.database-mysql-1", "module.nats"]
  source     = "./service-a"
}
module "service-2" {
  depends_on = ["module.network", "module.database-mysql-2", "module.nats"]
  source     = "./service-2"
}
module "service-3" {
  depends_on = ["module.network", "module.nats"]
  source     = "./service-3"
}
```

2. Set target file `export SHIPYARD_MODULES_HCL_FILE=/full/path/to/shipyard/modules/default.hcl`

### Run 
1. Install cli `go install github.com/Ishan27g/shpMod` 
2. From any directory run the command - `shpMod enable service-1`    
- disables all modules
- enables `service-1`, all modules `service-1` depends on, all modules those modules depend on ... and so on

```hcl
// Note -> Only the modules for service 1 are enabled
module "network" {
  source = "./network"
}
module "nats" {
  depends_on = ["module.network"]
  source     = "./nats"
}
module "database-mysql-1" {
  depends_on = ["module.network"]
  source     = "./mysql"
}
module "service-1" {
  depends_on = ["module.network", "module.database-mysql-1", "module.nats"]
  source     = "./service-a"
}
module "service-2" {
  disabled   = true
  depends_on = ["module.network", "module.database-mysql-2", "module.nats"]
  source     = "./service-2"
}
module "service-3" {
  disabled   = true
  depends_on = ["module.network", "module.nats"]
  source     = "./service-3"
}
module "database-mysql-2" {
  disabled   = true
  depends_on = ["module.network"]
  source     = "./mysql"
}
```

2. `shpMod run` & `shpMod destroy` are wrappers for `shipyard run` and `shipyard destroy` executed in the directory containing `SHIPYARD_MODULES_HCL_FILE`

> Autocomplete

Download [urfave/cli/autocomplete](https://github.com/urfave/cli/tree/main/autocomplete)
```shell
PROG=shpMOD
source path/to/downloaded/cli/autocomplete 
```

Commands `shpMod cmd arg` 
```shell
enable      e   -- enable modules `shpMod enable [tab] [tab] ...` | all `shpMod enable -all`
start       r   -- shipyard run                                                                                                                                                                              
stop        d   -- shipyard destroy                                                                                                                                                                             
setConfig   sc  -- set $HOME/.shpMod/cfg.hcl `shpMod setConfig ./config.hcl` |  clear $HOME/.shpMod/cfg.hcl `./shpMod setConfig -clear`                                                                                                                                                                 
showConfig  sh  -- print $HOME/.shpMod/cfg.hcl                                                                                                                                                               
```
3. Optional config. `config.hcl`
 - default modules that should be always enabled
 - names to skip when running `shpMod enable` shell completion  

```hcl
Config {
  enable = [  // optional
    "network",
  ]
  skipShellCompletion = [ // optional
    "nats",
    "network",
  ]
}
```