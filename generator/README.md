# Generators

## MenuItem

This directory is used to populate items for the daily quests menu GUI.  
Due to what appears to be BetonQuest limitations, the menu GUI is a tedious repetitive process. 
A standardization has been implemented across packages allowing for programatic creation of GUI menu entries.

### How to use

This automatically creates daily menu GUIs based on existing directories in the daily quest packages section.

Use the following in the base directory of the project to generate the GUI items.

```
go run menuitem.go
```

Enter the following values:  
- template type
- level

These values are directly tied to the file structure.

```QuestPackages/daily/<template_type>/<level>/```

This generator is to be used *after* the actual files are created.

The newly generated items should be found in the their respecitve template typed folder.
