# base.yml

This directory is used to populate items for the daily quests menu GUI.  
Due to what appears to be BetonQuest limitations, the menu GUI is a tedious repetitive process once a standardization has be implemented across packages. This standardization allows for programatically creating GUI menu entries.

To use the base.yml you fill in the values you need for a new quest.

Currently this creates a set of items that you can paste into the daily quests menu GUI. In the future this will automatically fill itself into the proper daily menu quest GUI, as well as create a new package for the new quest.

To generate the GUI items use the following in the base directory of the project.

```
go run main.go
```

The newly generated items should be found in the `generated` folder.