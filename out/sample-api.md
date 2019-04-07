# Name Search

**path**: addrBook/search

**sample uri**: http://namesearch.com/addrBook/search?fname=joe&lname=jackson&maxRec=389

* **maxRec**=99
* **fname**=  First name of person

* **lname**=  Last Name of person stored in addressbook.table.person.lname in primary oracle database. 
  Must Match legal name as shown on drivers license or passport.
  ```
  { 'person' : 
    {'lname' : 'myname' } 
  }
  ```

* **Type**= string
  **len**= 50


## Sample Output
```
  'person': {
    'lname': 'Jimbo',
    'fname': 'Jackson',
    'colors': {
       'car': 'red',
       'boat': 'blue',
       'house': 'cream',
       'hair': 'brown',
       'cat' : 'black and white spots'
    }
  } 
```

* **Colors** - List of colors for various items this person owns.  Used to help predict which color they will like when purchasing other items.

## Copyright


(C) Copyright Joseph Ellsworth Mar-2019

MIT License: https://opensource.org/licenses/MIT

Contact me on linkedin: https://www.linkedin.com/in/joe-ellsworth-68222/



