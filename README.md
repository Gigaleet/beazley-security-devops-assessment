# Beazley Security Tech Challenge DevOps

## Overview

Thank you for taking the time to apply to Beazley Security. Below is our technical challenge that we require individuals applying for our Software Engineering roles to complete. We recognize that time is valuable, and while you are free to spend as much time on this as you’d like, we have found that this will roughly take around an hour depending on experience.

The review of this challenge will focus on approach, style, reproducibility, and robustness from a software engineering perspective.

Please submit all code in a public repository and provide the link to the repository to our talent recruiter. After submission, we will complete an initial review and reach out to schedule a 2nd interview if submission is satisfactory. The results of this challenge will be discussed during the technical interview within our interview process which will be with technical members of our team.

---

## Task One

A 3-tier architecture is a common setup. Use a tool of your choosing/familiarity to create these.

**NFR’s:** There must always be 3 UP instances / containers within a cloud architecture of 3 zones (a/b/c). Please indicate within your solution how you would achieve this requirement.

---

## Task Two

Create a script that will query the meta data of an instance within a cloud provider of your choice and provide a JSON formatted output. The choice of language and implementation is up to you.

**NFR’s:** Your code allows for a particular data key to be retrieved individually, and unit testing is expected to be included.

---

## Task Three

You have a nested object, create a function that allows for the nested object and a key and returns the value of key. Any language besides Elixir is acceptable.

**Example Inputs**

| Object                 | Key    | Value |
|------------------------|--------|-------|
| {"a":{"b":{"c":"d"}}}  | a/b/c  | d     |
| {"x":{"y":{"z":"a"}}}  | x/y/z  | a     |

**NFR’s:** There must be Error Handling and unit testing in this challenge.
