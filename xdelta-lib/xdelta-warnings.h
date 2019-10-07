#ifndef _XDELTA_WARNINGS_837465732
#define _XDELTA_WARNINGS_837465732

#ifdef _MSC_VER
	#pragma warning( disable : 4100 )   // unused parameters
	#pragma warning( disable : 4244 )   // possible loss of data
    #pragma warning( disable : 4267 )   // possible loss of data
#endif

#ifdef __APPLE__
    #pragma clang diagnostic ignored "-Wunknown-warning-option"
    #pragma clang diagnostic ignored "-Wunused-parameter"
	#pragma clang diagnostic ignored "-Wreserved-user-defined-literal"
	#pragma clang diagnostic ignored "-Wc++11-extensions"
#endif

#ifdef __GNUC__
    #pragma GCC diagnostic ignored "-Wunknown-warning-option"
    #pragma GCC diagnostic ignored "-Wunused-parameter"
    #pragma GCC diagnostic ignored "-Wliteral-suffix"
#endif

#endif
