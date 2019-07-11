# Name Search

**path**: addrBook/search

**sample uri**: http://namesearch.com/addrBook/search?fname=joe&lname=jackson&maxRec=389

* **maxRec**= *(maxrec)* 99
* **fname**= {*db/person/fname} *FILE NOT FOUND data\db\person\fname.yml *
* **mname**= {*db/person/mname} *FILE NOT FOUND data\db\person\mname.yml *
* **lname**= {*db/person/lname#tech_desc} *FILE NOT FOUND data\db\person\lname.yml *
  * **Type**={*db/person/lname#type} *FILE NOT FOUND data\db\person\lname.yml *
  * ***len**={*db/person/lname#len} *FILE NOT FOUND data\db\person\lname.yml *
      * From Database desc 

## Sample Output
{*inc: share/person/example_person.md}

## Copyright  
{*inc: share/legal/copyright.txt}

