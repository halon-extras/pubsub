CMAKE_MINIMUM_REQUIRED(VERSION 3.5)

PROJECT("pubsub" VERSION 1.0.0)

# RPM/DEB and CENTOS8/ROCKY9/UBU2004/UBU2204/UBU2404
IF (EXISTS "${CMAKE_CURRENT_BINARY_DIR}/OSRELEASE.txt")
	FILE(STRINGS ${CMAKE_CURRENT_BINARY_DIR}/OSRELEASE.txt OSRELEASE)
ELSEIF (EXISTS "/OSRELEASE.txt")
	FILE(STRINGS /OSRELEASE.txt OSRELEASE)
ELSE()
	MESSAGE(FATAL_ERROR "Missing OSRELEASE")
ENDIF()

IF (OSRELEASE MATCHES "CENTOS8")
	SET(CPACK_GENERATOR "RPM")
ELSEIF (OSRELEASE MATCHES "ROCKY9")
	SET(CPACK_GENERATOR "RPM")
ELSEIF (OSRELEASE MATCHES "UBU2004")
	SET(CPACK_GENERATOR "DEB")
ELSEIF (OSRELEASE MATCHES "UBU2204")
	SET(CPACK_GENERATOR "DEB")
ELSEIF (OSRELEASE MATCHES "UBU2404")
	SET(CPACK_GENERATOR "DEB")
ELSE()
	MESSAGE(FATAL_ERROR "Invalid OSRELEASE")
ENDIF()

SET(CMAKE_BUILD_TYPE Release)

ADD_CUSTOM_TARGET(go ALL
	COMMAND ${CMAKE_COMMAND} -E env GOPATH=${GOPATH} /usr/local/go/bin/go build -buildmode c-shared -o ${CMAKE_CURRENT_BINARY_DIR}/pubsub.so
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

IF(EXISTS "/etc/os-release")
	EXECUTE_PROCESS(COMMAND "sed" "-ne" "s/^ID=\"\\?\\([a-z]\\+\\)\"\\?$/\\1/p" "/etc/os-release" OUTPUT_VARIABLE OS_RELEASE_ID OUTPUT_STRIP_TRAILING_WHITESPACE)
	EXECUTE_PROCESS(COMMAND "sed" "-ne" "s/^VERSION_ID=\"\\?\\([0-9\\.]\\+\\)\"\\?$/\\1/p" "/etc/os-release" OUTPUT_VARIABLE OS_RELEASE_VERSION_ID OUTPUT_STRIP_TRAILING_WHITESPACE)
	EXECUTE_PROCESS(COMMAND "uname" "-m" OUTPUT_VARIABLE OS_RELEASE_ARCH OUTPUT_STRIP_TRAILING_WHITESPACE)
	SET(SYSTEM_NAME "${OS_RELEASE_ID}-${OS_RELEASE_VERSION_ID}-${OS_RELEASE_ARCH}")
ENDIF()

INSTALL(FILES ${CMAKE_CURRENT_BINARY_DIR}/pubsub.so DESTINATION plugins)
INSTALL(FILES pubsub.schema.json DESTINATION share/plugins)
INSTALL(FILES pubsub.functions.json DESTINATION share/plugins)

IF(${CMAKE_VERSION} VERSION_GREATER_EQUAL "3.14.0") 
	FILE(CREATE_LINK "/opt/halon/plugins/pubsub.so" "${CMAKE_CURRENT_BINARY_DIR}/pubsub.so.link" SYMBOLIC)
ELSE()
	ADD_CUSTOM_TARGET(link_target ALL COMMAND ${CMAKE_COMMAND} -E create_symlink "/opt/halon/plugins/pubsub.so" "${CMAKE_CURRENT_BINARY_DIR}/pubsub.so.link")
ENDIF()

INSTALL(FILES "${CMAKE_CURRENT_BINARY_DIR}/pubsub.so.link" DESTINATION plugins/autoload/smtpd RENAME pubsub.so)
INSTALL(FILES "${CMAKE_CURRENT_BINARY_DIR}/pubsub.so.link" DESTINATION plugins/autoload/hsh RENAME pubsub.so)

INSTALL(FILES "LICENSE" DESTINATION share/doc/plugins/pubsub)
INSTALL(FILES "/licenses/LICENSE.go" DESTINATION share/doc/plugins/pubsub)
INSTALL(FILES "/licenses/LICENSE.google-cloud-go" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.httpsnoop" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.go-logr" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.s2a-go" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.uuid" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.enterprise-certificate-proxy" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.gax-go" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.opencensus-go" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.opentelemetry-go-instrumentation" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.google-api-go-client" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.go-genproto" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.grpc-go" DESTINATION share/doc/plugins/msal)
INSTALL(FILES "/licenses/LICENSE.protobuf-go" DESTINATION share/doc/plugins/msal)

SET(CMAKE_INSTALL_PREFIX "/opt/halon")
SET(CPACK_PACKAGING_INSTALL_PREFIX "/opt/halon")

IF (CPACK_GENERATOR MATCHES "DEB")
	SET(CPACK_DEBIAN_PACKAGE_DEPENDS "halon (>= 6.0)")
	SET(CPACK_DEBIAN_PACKAGE_CONTROL_EXTRA "${CMAKE_CURRENT_SOURCE_DIR}/postinst;${CMAKE_CURRENT_SOURCE_DIR}/postrm")
ENDIF()
IF (CPACK_GENERATOR MATCHES "RPM")
	SET(CPACK_RPM_PACKAGE_REQUIRES "halon >= 6.0")
	SET(CPACK_RPM_POST_INSTALL_SCRIPT_FILE "${CMAKE_CURRENT_SOURCE_DIR}/postinst")
	SET(CPACK_RPM_POST_UNINSTALL_SCRIPT_FILE "${CMAKE_CURRENT_SOURCE_DIR}/postrm")
ENDIF()

SET(CPACK_PACKAGE_NAME "halon-extras-pubsub")
SET(CPACK_PACKAGE_VERSION "${PROJECT_VERSION}")
SET(CPACK_PACKAGE_DESCRIPTION_SUMMARY "Halon extras pubsub")
SET(CPACK_PACKAGE_DESCRIPTION "Halon extras pubsub")
SET(CPACK_PACKAGE_CONTACT "Halon support@halon.io")
SET(CPACK_PACKAGE_FILE_NAME "${CPACK_PACKAGE_NAME}-${PROJECT_VERSION}-${SYSTEM_NAME}")

SET(CPACK_DEBIAN_PACKAGE_SECTION "mail")

INCLUDE(CPack)