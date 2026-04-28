---
trigger: always_on
---

Never forget [[api_standard_envelope]] ATD; Altering communication layer must trigger a warning for the user to approve.

Defaulting hides critical errors. Use them with caution. Crashing early is best (we are not in production mode)

Testing is meant to test production code: adding custom code to handle test should be made with great care and wisdom.

Don't use the script trigger_all_ci_tests.sh use trigger_one_ci_test.sh instead. The former is way too slow for pratical purpose and isn't easily stoppable, it erase too much informations. Note that run_all_unit_tests.sh exists as well and is quick enough to be abused. 