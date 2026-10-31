package main

import (
	"fmt"

	"github.com/eduardo-works/writeit/pkg/project"
	"github.com/eduardo-works/writeit/pkg/validator"
	"github.com/spf13/cobra"
)

func handleTitle(cmd *cobra.Command) (string, error) {
	title, err := cmd.Flags().GetString("title")
	if err != nil {
		return "", fmt.Errorf("title: %w", err)
	}

	return title, nil
}

func handlePositionalTitle(args []string, pos int) (string, error) {
	if pos < 0 || len(args) <= pos {
		return "", fmt.Errorf("title: %w: missing argument", project.ErrInvalidValue)
	}

	title, err := validator.ValidateNotEmpty(args[pos])
	if err != nil {
		return "", fmt.Errorf("title: %w", err)
	}

	return title, nil
}

func handleAuthor(cmd *cobra.Command) (string, error) {
	author, err := cmd.Flags().GetString("author")
	if err != nil {
		return "", fmt.Errorf("author: %w", err)
	}

	return author, nil
}

func handleEmail(cmd *cobra.Command) (string, error) {
	email, err := cmd.Flags().GetString("email")
	if err != nil {
		return "", fmt.Errorf("email: %w", err)
	}

	if email == "" {
		return "", nil
	}

	email, err = validator.ValidateEmail(email)
	if err != nil {
		return "", fmt.Errorf("email: %w", err)
	}

	return email, nil
}

func handlePhone(cmd *cobra.Command) (string, error) {
	phone, err := cmd.Flags().GetString("phone")
	if err != nil {
		return "", fmt.Errorf("phone: %w", err)
	}

	if phone == "" {
		return "", nil
	}

	phone, err = validator.ValidatePhone(phone)
	if err != nil {
		return "", fmt.Errorf("phone: %w", err)
	}

	return phone, nil
}

func handleTypeAndPart(cmd *cobra.Command) (project.ProjectType, int, error) {
	projectType, err := cmd.Flags().GetString("type")
	if err != nil {
		return project.Standalone, 0, fmt.Errorf("type: %w", err)
	}

	parts, err := cmd.Flags().GetInt("parts")
	if err != nil {
		return project.Standalone, 0, fmt.Errorf("parts: %w", err)
	}

	validatedType, validatedParts, err := validator.ValidateTypeAndParts(
		project.ProjectType(projectType),
		parts,
	)
	if err != nil {
		return project.Standalone, 0, fmt.Errorf("type/parts: %w", err)
	}

	return validatedType, validatedParts, nil
}
