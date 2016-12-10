#include <stdlib.h>
#include <stdio.h>

#if defined ( __APPLE__ )
#include <gl/glew.h>
#include <OpenCL/opencl.h>
#define GLFW_EXPOSE_NATIVE_COCOA
#define GLFW_EXPOSE_NATIVE_NSGL
#elif defined ( WIN32 )
#include <GL/glew.h>
#include <CL/cl.h>
#define GLFW_EXPOSE_NATIVE_WIN32
#define GLFW_EXPOSE_NATIVE_WGL
#else
#include <GL/glew.h>
#include <CL/cl.h>
#define GLFW_EXPOSE_NATIVE_X11
#define GLFW_EXPOSE_NATIVE_GLX
#endif

#include <GLFW/glfw3.h>
#include <GLFW/glfw3native.h>

#include "imgui.h"
#include "imgui_impl_glfw_gl3.h"

#ifdef __cplusplus
extern "C" {
#endif

  GLFWwindow *window;

  static void error_callback(int error, const char* description) {
    fprintf(stderr, "Error (%d): %s\n", error, description);
    exit(EXIT_FAILURE);
  }

  void graphview_init(int width, int height, int fullscreen) {
    if (!glfwInit()) {
      exit(EXIT_FAILURE);
    }

    glfwSetErrorCallback(error_callback);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 3);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 3);
    glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, 1);
    glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);
    glfwWindowHint(GLFW_VISIBLE, 1);

    const GLFWvidmode *vidMode = glfwGetVideoMode(glfwGetPrimaryMonitor());

    if (fullscreen == 1) {
      // ignore requested width and height, we're going fullscreen!
      window = glfwCreateWindow(vidMode->width, vidMode->height, "Konig", glfwGetPrimaryMonitor(), NULL);
    } else {
      window = glfwCreateWindow(width, height, "Konig", NULL, NULL);
    }

    glfwMakeContextCurrent(window);
    glewExperimental = GL_TRUE;
    glewInit();

    // init imgui
    ImGui_ImplGlfwGL3_Init(window, true);
    ImGuiIO& io = ImGui::GetIO();
    //io.WantCaptureMouse = true;
    io.MouseDrawCursor = true;
    io.IniFilename = NULL;

    // enable v-sync
    glfwSwapInterval(1);

    // set immutable gl stuff
    glClearColor(0.7, 0.7, 0.7, 1.0);
  }

  ImVec4 clear_color = ImColor(114, 144, 154);

  int graphview_update(double *dt) {
    // draw our stuff
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    {}

    // draw ui
    ImGui_ImplGlfwGL3_NewFrame();

    bool show_test_window = true;

    // 1. Show a simple window
    // Tip: if we don't call ImGui::Begin()/ImGui::End() the widgets appears in a window automatically called "Debug"
    {
      static float f = 0.0f;
      ImGui::Text("Hello, world!");
      ImGui::SliderFloat("float", &f, 0.0f, 1.0f);
      ImGui::ColorEdit3("clear color", (float*)&clear_color);
      if (ImGui::Button("Test Window")) show_test_window ^= 1;
      ImGui::Text("Application average %.3f ms/frame (%.1f FPS)", 1000.0f / ImGui::GetIO().Framerate, ImGui::GetIO().Framerate);
    }

    ImGui::Render();
    glfwSwapBuffers(window);

    *dt = 1.0/ImGui::GetIO().Framerate;

    // get events
    glfwPollEvents();
    return glfwWindowShouldClose(window);
  }

  void graphview_shutdown() {
    ImGui_ImplGlfwGL3_Shutdown();
    glfwDestroyWindow(window);
    glfwTerminate();
  }

#ifdef __cplusplus
}
#endif
