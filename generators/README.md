# Generators

## MenuItem

This directory is used to populate items for the daily quests menu GUI.  
Due to what appears to be BetonQuest limitations, the menu GUI is a tedious repetitive process once a standardization has be implemented across packages. This standardization allows for programatically creating GUI menu entries.

### How to use

Currently this creates a set of items that you can paste into the daily quests menu GUI. In the future this will automatically fill itself into the proper daily menu quest GUI.

To generate the GUI items use the following in the base directory of the project.

```
go run menuitem.go
```

Enter the following three values:  
- template type
- tier
- filename

These three types are the file directly tied to the file structure.

QuestPackages/daily/<template_type>/<tier>/<filename>

This generator is to be used *after* the actual files are created.

The newly generated items should be found in the `generated` folder.