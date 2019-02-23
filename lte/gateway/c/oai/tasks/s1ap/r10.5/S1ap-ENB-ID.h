/*
 * Generated by asn1c-0.9.24 (http://lionet.info/asn1c)
 * From ASN.1 module "S1AP-IEs"
 * 	found in "/home/vagrant/oai/s1ap/messages/asn1/R10.5/S1AP-IEs.asn"
 * 	`asn1c -gen-PER`
 */

#ifndef _S1ap_ENB_ID_H_
#define _S1ap_ENB_ID_H_

#include <asn_application.h>

/* Including external dependencies */
#include <BIT_STRING.h>
#include <constr_CHOICE.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Dependencies */
typedef enum S1ap_ENB_ID_PR {
  S1ap_ENB_ID_PR_NOTHING, /* No components present */
  S1ap_ENB_ID_PR_macroENB_ID,
  S1ap_ENB_ID_PR_homeENB_ID,
  /* Extensions may appear below */

} S1ap_ENB_ID_PR;

/* S1ap-ENB-ID */
typedef struct S1ap_ENB_ID {
  S1ap_ENB_ID_PR present;
  union S1ap_ENB_ID_u {
    BIT_STRING_t macroENB_ID;
    BIT_STRING_t homeENB_ID;
    /*
		 * This type is extensible,
		 * possible extensions are below.
		 */
  } choice;

  /* Context for parsing across buffer boundaries */
  asn_struct_ctx_t _asn_ctx;
} S1ap_ENB_ID_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_S1ap_ENB_ID;

#ifdef __cplusplus
}
#endif

#endif /* _S1ap_ENB_ID_H_ */
#include <asn_internal.h>
