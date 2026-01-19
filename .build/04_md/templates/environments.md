# {{ name|upper }}

***Tier {{ tier }} {{ type }}***  
*{{ description }}*  
**Impulses:** {{ impulses }}

> **Difficulty:** {{ difficulty }}  
> **Potential Adversaries:** {{ potential_adversaries }}

{% set questions = [
  feature_question_1, feature_question_2, feature_question_3,
  feature_question_4, feature_question_5, feature_question_6,
  feature_question_7, feature_question_8, feature_question_9,
  feature_question_10, feature_question_11, feature_question_12,
  feature_question_13, feature_question_14, feature_question_15
] %}
{% set questions = questions | select | list %}
{% if questions %}
## FEATURE QUESTIONS
{% for question in questions -%}
- {{ question }}
{% endfor %}
{% endif %}

## FEATURES
{% for feat in feats %}
***{{ feat.name }}:*** {{ feat.text }}
{% endfor %}
