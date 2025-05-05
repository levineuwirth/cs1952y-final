/*
 *
 * FILE: analyze_simd.go
 * LATEST: 10:08 05 May 2025
 * DESC: find percentage of a dumped amd64 object file's instructions that are SIMD instructions
 * AUTHOR: Levi Neuwirth <ln@levineuwirth.org>
 *
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var total int
var simd int

// Since Go doesn't have a hashset, we will use a hashmap and ignore the Value...
var simdInstr map[string]bool
var digits []string

func main() {
	if len(os.Args) < 1 {
		log.Fatal("Usage: ./analyze_simd <path to .txt from objdump>")
	}

	objDumpRaw, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully opened object dump. Investigating...")
	}
	defer objDumpRaw.Close()

	initDigits()
	initSimdInstructions()

	// This regex magic will get us the instructions from an extracted objdump line.
	instrRegex := regexp.MustCompile(`\b([a-z]{2,6}[a-z]*)\b`)

	scanner := bufio.NewScanner(objDumpRaw)
	for scanner.Scan() {
		localLine := scanner.Text()
		localLineSplit := strings.Fields(localLine)

		if len(localLineSplit) < 2 || !strings.Contains(localLineSplit[0], ":") {
			continue
		}

		matches := instrRegex.FindAllString(localLine, -1)
		if len(matches) == 0 {
			continue
		}

		instr := matches[0]
		log.Println(instr)
		if simdInstr[instr] {
			simd++
		}
		total++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The result is:\n%d SIMD instructions\n%d Total instructions\n", simd, total)
}

func initSimdInstructions() {
	simdInstr = map[string]bool{
		// MMX Instructions
		"packsswb": true, "packssdw": true,
		"packuswb": true, "paddb": true, "paddw": true, "paddd": true,
		"paddsb": true, "paddsw": true, "paddusb": true, "paddusw": true,
		"pand": true, "pandn": true, "pcmpeqb": true, "pcmpeqw": true,
		"pcmpeqd": true, "pcmpgtb": true, "pcmpgtw": true, "pcmpgtd": true,
		"pmaddwd": true, "pmulhw": true, "pmullw": true, "por": true,
		"psllw": true, "pslld": true, "psllq": true, "psraw": true,
		"psrad": true, "psrlw": true, "psrld": true, "psrlq": true,
		"psubb": true, "psubw": true, "psubd": true, "psubsb": true,
		"psubsw": true, "psubusb": true, "psubusw": true, "punpckhbw": true,
		"punpckhwd": true, "punpckhdq": true, "punpcklbw": true, "punpcklwd": true,
		"punpckldq": true, "pxor": true,

		// SSE Instructions
		"addps": true, "addss": true, "andps": true, "andnps": true,
		"cmpeqps": true, "cmpeqss": true, "cmpgeps": true, "cmpgess": true,
		"cmpgtps": true, "cmpgtss": true, "cmpleps": true, "cmpless": true,
		"cmpltps": true, "cmpltss": true, "cmpneqps": true, "cmpneqss": true,
		"cmpngeps": true, "cmpngess": true, "cmpngtps": true, "cmpngtss": true,
		"cmpnleps": true, "cmpnless": true, "cmpnltps": true, "cmpnltss": true,
		"cmpordps": true, "cmpordss": true, "cmpunordps": true, "cmpunordss": true,
		"divps": true, "divss": true, "maxps": true, "maxss": true,
		"minps": true, "minss": true, "movaps": true, "movss": true,
		"movups": true, "mulps": true, "mulss": true, "rcpps": true,
		"rcpss": true, "rsqrtps": true, "rsqrtss": true, "sqrtps": true,
		"sqrtss": true, "subps": true, "subss": true, "xorps": true,

		// SSE2 Instructions
		"addpd": true, "addsd": true, "andpd": true, "andnpd": true,
		"cmpeqpd": true, "cmpeqsd": true, "cmpgepd": true, "cmpgesd": true,
		"cmpgtpd": true, "cmpgtsd": true, "cmplepd": true, "cmplesd": true,
		"cmpltpd": true, "cmpltsd": true, "cmpneqpd": true, "cmpneqsd": true,
		"cmpngepd": true, "cmpngesd": true, "cmpngtpd": true, "cmpngtsd": true,
		"cmpnlepd": true, "cmpnlesd": true, "cmpnltpd": true, "cmpnltsd": true,
		"cmpordpd": true, "cmpordsd": true, "cmpunordpd": true, "cmpunordsd": true,
		"divpd": true, "divsd": true, "maxpd": true, "maxsd": true,
		"minpd": true, "minsd": true, "movapd": true, "movsd": true,
		"movupd": true, "mulpd": true, "mulsd": true, "sqrtpd": true,
		"subpd": true, "subsd": true, "xorpd": true,

		// SSE3 Instructions
		"addsubpd": true, "addsubps": true, "haddpd": true, "haddps": true,
		"hsubpd": true, "hsubps": true, "lddqu": true, "monitor": true,
		"mwait": true, "movddup": true, "movshdup": true, "movsldup": true,

		// SSSE3 Instructions
		"pshufb": true, "phaddw": true, "phaddd": true, "phaddsw": true,
		"pmaddubsw": true, "phsubw": true, "phsubd": true, "phsubsw": true,
		"psignb": true, "psignw": true, "psignd": true, "pmulhrsw": true,
		"palignr": true,

		// SSE4.1 Instructions
		"blendpd": true, "blendps": true, "blendvpd": true, "blendvps": true,
		"dppd": true, "dpps": true, "extractps": true, "insertps": true,
		"movntdqa": true, "mpsadbw": true, "packusdw": true, "pblendvb": true,
		"pblendw": true, "pcmpeqq": true, "pextrb": true, "pextrd": true,
		"pextrq": true, "phminposuw": true, "pinsrb": true, "pinsrd": true,
		"pinsrq": true, "pmuldq": true, "pmulld": true, "ptest": true,
		"roundpd": true, "roundps": true, "roundsd": true, "roundss": true,

		// SSE4.2 Instructions
		"pcmpestri": true, "pcmpestrm": true, "pcmpistri": true, "pcmpistrm": true,
		"crc32": true, "popcnt": true,

		// AVX Instructions
		"vaddpd": true, "vaddps": true, "vaddsd": true, "vaddss": true,
		"vandpd": true, "vandps": true, "vandnpd": true, "vandnps": true,
		"vdivpd": true, "vdivps": true, "vdivsd": true, "vdivss": true,
		"vmaxpd": true, "vmaxps": true, "vmaxsd": true, "vmaxss": true,
		"vminpd": true, "vminps": true, "vminsd": true, "vminss": true,
		"vmulpd": true, "vmulps": true, "vmulsd": true, "vmulss": true,
		"vorpd": true, "vorps": true, "vsqrtpd": true, "vsqrtps": true,
		"vsqrtsd": true, "vsqrtss": true, "vsubpd": true, "vsubps": true,
		"vsubsd": true, "vsubss": true, "vxorpd": true, "vxorps": true,

		// AVX2 Instructions
		"vpabsb": true, "vpabsw": true, "vpabsd": true, "vpaddb": true,
		"vpaddw": true, "vpaddd": true, "vpaddq": true, "vpaddsb": true,
		"vpaddsw": true, "vpaddusb": true, "vpaddusw": true, "vpalignr": true,
		"vpand": true, "vpandn": true, "vpavgb": true, "vpavgw": true,
		"vpblendd": true, "vpcmpeqb": true, "vpcmpeqw": true, "vpcmpeqd": true,
		"vpcmpeqq": true, "vpcmpgtb": true, "vpcmpgtw": true, "vpcmpgtd": true,

		// AVX512 not included since Kyber does not use it.
	}
}

func initDigits() {
	digits = make([]string, 0)
	digits = append(digits, "0")
	digits = append(digits, "1")
	digits = append(digits, "2")
	digits = append(digits, "3")
	digits = append(digits, "4")
	digits = append(digits, "5")
	digits = append(digits, "6")
	digits = append(digits, "7")
	digits = append(digits, "8")
	digits = append(digits, "9")
}
