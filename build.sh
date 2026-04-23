#!/bin/bash

# MDCLI - Markdown CLI 工具编译脚本
# 支持 AMD64 和 ARM64 架构

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目信息
APP_NAME="mdcli"
VERSION="2.0"
BUILD_DIR="build"
SOURCE_FILE="."

# 打印带颜色的信息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
MDCLI 工具编译脚本

用法: $0 [选项]

选项:
    -h, --help              显示此帮助信息
    -a, --all               编译所有平台和架构
    -l, --linux             编译 Linux 平台（amd64 和 arm64）
    -m, --macos             编译 macOS 平台（amd64 和 arm64）
    -w, --windows           编译 Windows 平台（amd64 和 arm64）
    -c, --clean             清理构建目录
    --linux-amd64           仅编译 Linux AMD64
    --linux-arm64           仅编译 Linux ARM64
    --darwin-amd64          仅编译 macOS AMD64 (Intel)
    --darwin-arm64          仅编译 macOS ARM64 (Apple Silicon)
    --windows-amd64         仅编译 Windows AMD64
    --windows-arm64         仅编译 Windows ARM64

示例:
    $0 --all                # 编译所有平台
    $0 --linux              # 编译 Linux 平台
    $0 --darwin-arm64       # 仅编译 macOS ARM64
    $0 --clean              # 清理构建目录

EOF
}

# 清理构建目录
clean_build() {
    print_info "清理构建目录..."
    if [ -d "$BUILD_DIR" ]; then
        rm -rf "$BUILD_DIR"
        print_success "构建目录已清理"
    else
        print_warning "构建目录不存在"
    fi
}

# 创建构建目录
prepare_build_dir() {
    if [ ! -d "$BUILD_DIR" ]; then
        mkdir -p "$BUILD_DIR"
        print_info "创建构建目录: $BUILD_DIR"
    fi
}

# 编译函数
build_binary() {
    local os=$1
    local arch=$2
    local output_name="${APP_NAME}"
    
    # Windows 需要 .exe 后缀
    if [ "$os" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    local output_path="${BUILD_DIR}/${APP_NAME}-${os}-${arch}"
    mkdir -p "$output_path"
    local output_file="${output_path}/${output_name}"
    
    print_info "编译 ${os}/${arch}..."
    
    # 设置环境变量并编译
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -mod=vendor \
        -ldflags="-s -w" \
        -o "$output_file" \
        "$SOURCE_FILE"
    
    if [ $? -eq 0 ]; then
        local size=$(du -h "$output_file" | cut -f1)
        print_success "编译成功: $output_file (大小: $size)"
        
        # 创建压缩包
        cd "$BUILD_DIR"
        local archive_name="${APP_NAME}-${os}-${arch}-v${VERSION}"
        if [ "$os" = "windows" ]; then
            zip -q -r "${archive_name}.zip" "${APP_NAME}-${os}-${arch}"
            print_success "已创建压缩包: ${archive_name}.zip"
        else
            tar -czf "${archive_name}.tar.gz" "${APP_NAME}-${os}-${arch}"
            print_success "已创建压缩包: ${archive_name}.tar.gz"
        fi
        cd ..
    else
        print_error "编译失败: ${os}/${arch}"
        return 1
    fi
}

# 编译所有平台
build_all() {
    print_info "开始编译所有平台..."
    echo ""
    
    # Linux
    build_binary "linux" "amd64"
    build_binary "linux" "arm64"
    
    # macOS
    build_binary "darwin" "amd64"
    build_binary "darwin" "arm64"
    
    # Windows
    build_binary "windows" "amd64"
    build_binary "windows" "arm64"
    
    echo ""
    print_success "所有平台编译完成！"
    show_build_summary
}

# 编译 Linux 平台
build_linux() {
    print_info "编译 Linux 平台..."
    build_binary "linux" "amd64"
    build_binary "linux" "arm64"
    show_build_summary
}

# 编译 macOS 平台
build_macos() {
    print_info "编译 macOS 平台..."
    build_binary "darwin" "amd64"
    build_binary "darwin" "arm64"
    show_build_summary
}

# 编译 Windows 平台
build_windows() {
    print_info "编译 Windows 平台..."
    build_binary "windows" "amd64"
    build_binary "windows" "arm64"
    show_build_summary
}

# 显示编译摘要
show_build_summary() {
    echo ""
    print_info "编译摘要:"
    echo "----------------------------------------"
    if [ -d "$BUILD_DIR" ]; then
        ls -lh "$BUILD_DIR" | grep -E "\.(tar\.gz|zip)$" | awk '{print "  " $9 " - " $5}'
    fi
    echo "----------------------------------------"
    print_info "构建目录: $BUILD_DIR"
}

# 检查 Go 环境
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "未找到 Go 环境，请先安装 Go"
        exit 1
    fi
    
    local go_version=$(go version | awk '{print $3}')
    print_info "Go 版本: $go_version"
}

# 主函数
main() {
    echo ""
    print_info "========================================"
    print_info "  MDCLI 工具编译脚本 v${VERSION}"
    print_info "========================================"
    echo ""
    
    # 检查 Go 环境
    check_go
    
    # 如果没有参数，显示帮助
    if [ $# -eq 0 ]; then
        show_help
        exit 0
    fi
    
    # 准备构建目录
    prepare_build_dir
    
    # 解析参数
    case "$1" in
        -h|--help)
            show_help
            ;;
        -c|--clean)
            clean_build
            ;;
        -a|--all)
            build_all
            ;;
        -l|--linux)
            build_linux
            ;;
        -m|--macos)
            build_macos
            ;;
        -w|--windows)
            build_windows
            ;;
        --linux-amd64)
            build_binary "linux" "amd64"
            show_build_summary
            ;;
        --linux-arm64)
            build_binary "linux" "arm64"
            show_build_summary
            ;;
        --darwin-amd64)
            build_binary "darwin" "amd64"
            show_build_summary
            ;;
        --darwin-arm64)
            build_binary "darwin" "arm64"
            show_build_summary
            ;;
        --windows-amd64)
            build_binary "windows" "amd64"
            show_build_summary
            ;;
        --windows-arm64)
            build_binary "windows" "arm64"
            show_build_summary
            ;;
        *)
            print_error "未知选项: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@"

