# Name Search

**path**: addrBook/search

**sample uri**: http://namesearch.com/addrBook/search?fname=joe&lname=jackson&maxRec=389

* **maxRec**=99
* **fname**=  First name of person
type: string
len: 50
* **lname**=  Last Name of person stored in addressbook.table.person.lname in primary oracle database. 
  Must Match legal name as shown on drivers license or passport.
  ```
  { 'person' : 
    {'lname' : 'myname' } 
  }
  ```
type: string
len: 50
* **Type**= string
len: 50  **len**= 50




(C) Copyright Joseph Ellsworth Mar-2019

MIT License: https://opensource.org/licenses/MIT

Contact me on linkedin: https://www.linkedin.com/in/joe-ellsworth-68222/

