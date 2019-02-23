/*
 * Generated by asn1c-0.9.24 (http://lionet.info/asn1c)
 * From ASN.1 module "S1AP-IEs"
 * 	found in "/home/vagrant/oai/s1ap/messages/asn1/R10.5/S1AP-IEs.asn"
 * 	`asn1c -gen-PER`
 */

#ifndef _S1ap_ENBname_H_
#define _S1ap_ENBname_H_

#include <asn_application.h>

/* Including external dependencies */
#include <PrintableString.h>

#ifdef __cplusplus
extern "C" {
#endif

/* S1ap-ENBname */
typedef PrintableString_t S1ap_ENBname_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_S1ap_ENBname;
asn_struct_free_f S1ap_ENBname_free;
asn_struct_print_f S1ap_ENBname_print;
asn_constr_check_f S1ap_ENBname_constraint;
ber_type_decoder_f S1ap_ENBname_decode_ber;
der_type_encoder_f S1ap_ENBname_encode_der;
xer_type_decoder_f S1ap_ENBname_decode_xer;
xer_type_encoder_f S1ap_ENBname_encode_xer;
per_type_decoder_f S1ap_ENBname_decode_uper;
per_type_encoder_f S1ap_ENBname_encode_uper;
per_type_decoder_f S1ap_ENBname_decode_aper;
per_type_encoder_f S1ap_ENBname_encode_aper;
type_compare_f S1ap_ENBname_compare;

#ifdef __cplusplus
}
#endif

#endif /* _S1ap_ENBname_H_ */
#include <asn_internal.h>
