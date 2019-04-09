# Name Search

**path**: addrBook/search

**sample uri**: http://namesearch.com/addrBook/search?fname=joe&lname=jackson&maxRec=389

* **maxRec**= *(maxrec)* 99
* **fname**=  *(db/person/fname)*  First name of person

* **lname**=  *(db/person/lname#tech_desc)*  Last Name of person stored in addressbook.table.person.lname in primary oracle database. 
  Must Match legal name as shown on drivers license or passport. 
  ```
  { 'person' : 
    {'lname' : 'myname' } 
  }
  ```

* **Type**= *(db/person/lname#type)*  string
  **len**= *(db/person/lname#len)*  50


## Sample Output
 *(inc: share/person/example_person.md)*  
```
  'person': {
    'lname': 'Jimbo',
    'fname': 'Jackson',
    'colors': {
       'car': 'red',
       'boat': 'blue',
       'house': 'cream',
       'hair': 'dark brown',
       'cat' : 'black and white spots'
    }
  } 
```

* **Colors** - List of colors for various items this person owns.  Used to help predict which color they will like when purchasing other items.

## Copyright
 *(inc: share/legal/copyright.txt)*  


(C) Copyright Joseph Ellsworth Mar-2019

MIT License: https://opensource.org/licenses/MIT

Contact me on linkedin: https://www.linkedin.com/in/joe-ellsworth-68222/ if you want a great programmer.



