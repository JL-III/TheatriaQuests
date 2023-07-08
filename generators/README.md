# Generators

## MenuItem

This directory is used to populate items for the daily quests menu GUI.  
Due to what appears to be BetonQuest limitations, the menu GUI is a tedious repetitive process once a standardization has be implemented across packages. This standardization allows for programatically creating GUI menu entries.

### How to use

This automatically creates daily menu guis based on existing directories in the daily quest packages section.

To generate the GUI items use the following in the base directory of the project.

```
go run menuitem.go
```

Enter the following values:  
- template type
- tier

These values are directly tied to the file structure.

```QuestPackages/daily/<template_type>/<tier>/```

This generator is to be used *after* the actual files are created.

The newly generated items should be found in the their respecitve template typed folder.
