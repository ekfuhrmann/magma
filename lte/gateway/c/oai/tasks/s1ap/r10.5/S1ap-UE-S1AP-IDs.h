/*
 * Generated by asn1c-0.9.24 (http://lionet.info/asn1c)
 * From ASN.1 module "S1AP-IEs"
 * 	found in "/home/vagrant/oai/s1ap/messages/asn1/R10.5/S1AP-IEs.asn"
 * 	`asn1c -gen-PER`
 */

#ifndef _S1ap_UE_S1AP_IDs_H_
#define _S1ap_UE_S1AP_IDs_H_

#include <asn_application.h>

/* Including external dependencies */
#include "S1ap-UE-S1AP-ID-pair.h"
#include "S1ap-MME-UE-S1AP-ID.h"
#include <constr_CHOICE.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Dependencies */
typedef enum S1ap_UE_S1AP_IDs_PR {
  S1ap_UE_S1AP_IDs_PR_NOTHING, /* No components present */
  S1ap_UE_S1AP_IDs_PR_uE_S1AP_ID_pair,
  S1ap_UE_S1AP_IDs_PR_mME_UE_S1AP_ID,
  /* Extensions may appear below */

} S1ap_UE_S1AP_IDs_PR;

/* S1ap-UE-S1AP-IDs */
typedef struct S1ap_UE_S1AP_IDs {
  S1ap_UE_S1AP_IDs_PR present;
  union S1ap_UE_S1AP_IDs_u {
    S1ap_UE_S1AP_ID_pair_t uE_S1AP_ID_pair;
    S1ap_MME_UE_S1AP_ID_t mME_UE_S1AP_ID;
    /*
		 * This type is extensible,
		 * possible extensions are below.
		 */
  } choice;

  /* Context for parsing across buffer boundaries */
  asn_struct_ctx_t _asn_ctx;
} S1ap_UE_S1AP_IDs_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_S1ap_UE_S1AP_IDs;

#ifdef __cplusplus
}
#endif

#endif /* _S1ap_UE_S1AP_IDs_H_ */
#include <asn_internal.h>
