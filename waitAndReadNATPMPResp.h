#ifndef _WAITANDREADNATPMPRESP
#define _WAITANDREADNATPMPRESP

#include "natpmp.h"

typedef struct {
	in_addr_t addr;
	uint16_t privateport;
	uint16_t mappedpublicport;
	uint32_t lifetime;
} resp_s;

int waitAndReadNATPMPResp(natpmp_t *natpmp, resp_s *resp);

#endif